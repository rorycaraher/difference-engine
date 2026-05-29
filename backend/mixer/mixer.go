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

type Config struct {
	// StemsDir and OutputDir are base directories; the track slug is appended as a subdirectory.
	StemsDir          string
	OutputDir         string
	R2AccountID       string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2StemsBucket     string
}

type Mixer struct {
	cfg Config
}

func New(cfg Config) *Mixer {
	return &Mixer{cfg: cfg}
}

// GetStems returns local file paths or R2 presigned URLs for the given track and stem names.
// Local path:  StemsDir/track/stem.mp3
// R2 key:      track/stem.mp3
func (m *Mixer) GetStems(track string, selectedStems []string) ([]string, error) {
	if m.cfg.R2StemsBucket != "" {
		return m.presignedURLs(track, selectedStems)
	}
	paths := make([]string, len(selectedStems))
	for i, stem := range selectedStems {
		paths[i] = filepath.Join(m.cfg.StemsDir, track, stem+".mp3")
	}
	return paths, nil
}

func (m *Mixer) presignedURLs(track string, stems []string) ([]string, error) {
	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			m.cfg.R2AccessKeyID, m.cfg.R2SecretAccessKey, "",
		)),
		awsconfig.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", m.cfg.R2AccountID))
		o.UsePathStyle = true
	})
	presigner := s3.NewPresignClient(client)

	urls := make([]string, len(stems))
	for i, stem := range stems {
		resp, err := presigner.PresignGetObject(context.Background(), &s3.GetObjectInput{
			Bucket: aws.String(m.cfg.R2StemsBucket),
			Key:    aws.String(track + "/" + stem + ".mp3"),
		}, func(o *s3.PresignOptions) {
			o.Expires = 5 * time.Minute
		})
		if err != nil {
			return nil, fmt.Errorf("presign %s/%s: %w", track, stem, err)
		}
		urls[i] = resp.URL
	}
	return urls, nil
}

// CreateMixdown mixes inputFiles at the given volumes via ffmpeg.
// Output is written to OutputDir/track/output-<timestamp>.mp3.
// If volumes is nil, random values in [0.5, 1.0] are used.
func (m *Mixer) CreateMixdown(track string, inputFiles []string, volumes []float64) (string, error) {
	if volumes == nil {
		volumes = make([]float64, len(inputFiles))
		for i := range volumes {
			v := rand.Float64()*0.5 + 0.5
			volumes[i] = math.Round(v*1000) / 1000
		}
	}

	outDir := filepath.Join(m.cfg.OutputDir, track)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return "", fmt.Errorf("create output dir: %w", err)
	}

	timestamp := time.Now().Format("20060102150405")
	outputPath := filepath.Join(outDir, fmt.Sprintf("output-%s.mp3", timestamp))

	cmd := exec.Command("ffmpeg", buildFFmpegArgs(inputFiles, volumes, outputPath)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg: %w", err)
	}

	return outputPath, nil
}

// buildFFmpegArgs constructs the equivalent of:
//
//	streams = [ffmpeg.input(f).filter("volume", vol) for f, vol in zip(inputs, volumes)]
//	ffmpeg.filter(streams, "amix", inputs=N, duration="longest").output(out, acodec="libmp3lame", audio_bitrate="192k")
func buildFFmpegArgs(inputFiles []string, volumes []float64, outputPath string) []string {
	args := make([]string, 0, len(inputFiles)*2+6)
	for _, f := range inputFiles {
		args = append(args, "-i", f)
	}

	parts := make([]string, 0, len(inputFiles)+1)
	var mixInputs strings.Builder
	for i, vol := range volumes {
		label := fmt.Sprintf("v%d", i)
		parts = append(parts, fmt.Sprintf("[%d:a]volume=%.3f[%s]", i, vol, label))
		fmt.Fprintf(&mixInputs, "[%s]", label)
	}
	parts = append(parts, fmt.Sprintf(
		"%samix=inputs=%d:duration=longest[mixed]",
		mixInputs.String(), len(inputFiles),
	))
	parts = append(parts, "[mixed]loudnorm=I=-14:TP=-1:LRA=11[mix]")

	return append(args,
		"-filter_complex", strings.Join(parts, ";"),
		"-map", "[mix]",
		"-acodec", "libmp3lame",
		"-b:a", "192k",
		outputPath,
	)
}
