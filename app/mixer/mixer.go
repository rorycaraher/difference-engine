package mixer

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Mixer struct{}

// GetStems returns local file paths, or R2 presigned URLs when R2_STEMS_BUCKET is set.
func (m *Mixer) GetStems(directory string, selectedStems []string) ([]string, error) {
	if R2StemsBucket != "" {
		return presignedURLs(selectedStems)
	}
	paths := make([]string, len(selectedStems))
	for i, stem := range selectedStems {
		paths[i] = filepath.Join(directory, stem+".mp3")
	}
	return paths, nil
}

func presignedURLs(stems []string) ([]string, error) {
	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			R2AccessKeyID, R2SecretAccessKey, "",
		)),
		awsconfig.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", R2AccountID))
		o.UsePathStyle = true
	})
	presigner := s3.NewPresignClient(client)

	urls := make([]string, len(stems))
	for i, stem := range stems {
		resp, err := presigner.PresignGetObject(context.Background(), &s3.GetObjectInput{
			Bucket: aws.String(R2StemsBucket),
			Key:    aws.String(stem + ".mp3"),
		}, func(o *s3.PresignOptions) {
			o.Expires = 5 * time.Minute
		})
		if err != nil {
			return nil, fmt.Errorf("presign %s: %w", stem, err)
		}
		urls[i] = resp.URL
	}
	return urls, nil
}

// CreateMixdown mixes inputFiles at the given volumes via ffmpeg and returns the output path.
// If volumes is nil, random values in [0.5, 1.0] are used.
func (m *Mixer) CreateMixdown(inputFiles []string, volumes []float64) (string, error) {
	if volumes == nil {
		volumes = make([]float64, len(inputFiles))
		for i := range volumes {
			v := rand.Float64()*0.5 + 0.5
			volumes[i] = math.Round(v*1000) / 1000
		}
	}

	if err := os.MkdirAll(OutputDir, 0755); err != nil {
		return "", fmt.Errorf("create output dir: %w", err)
	}

	timestamp := time.Now().Format("20060102150405")
	outputPath := filepath.Join(OutputDir, fmt.Sprintf("output-%s.mp3", timestamp))

	args := buildFFmpegArgs(inputFiles, volumes, outputPath)
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg: %w", err)
	}

	return outputPath, nil
}

// buildFFmpegArgs constructs the ffmpeg argument list equivalent to the Python ffmpeg-python call:
//
//	streams = [ffmpeg.input(f).filter("volume", vol) for f, vol in zip(inputs, volumes)]
//	ffmpeg.filter(streams, "amix", inputs=N, duration="longest").output(out, acodec="libmp3lame", audio_bitrate="192k")
func buildFFmpegArgs(inputFiles []string, volumes []float64, outputPath string) []string {
	args := make([]string, 0, len(inputFiles)*2+6)
	for _, f := range inputFiles {
		args = append(args, "-i", f)
	}

	// Build filter_complex: volume filter per stream, then amix
	parts := make([]string, 0, len(inputFiles)+1)
	var mixInputs strings.Builder
	for i, vol := range volumes {
		label := fmt.Sprintf("v%d", i)
		parts = append(parts, fmt.Sprintf("[%d:a]volume=%.3f[%s]", i, vol, label))
		fmt.Fprintf(&mixInputs, "[%s]", label)
	}
	parts = append(parts, fmt.Sprintf(
		"%samix=inputs=%d:duration=longest[mix]",
		mixInputs.String(), len(inputFiles),
	))

	args = append(args,
		"-filter_complex", strings.Join(parts, ";"),
		"-map", "[mix]",
		"-acodec", "libmp3lame",
		"-b:a", "192k",
		outputPath,
	)
	return args
}
