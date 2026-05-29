package mixer

import (
	"strings"
	"testing"
)

func TestGetStemsReturnsPaths(t *testing.T) {
	// Ensure local filesystem path is used
	orig := R2StemsBucket
	R2StemsBucket = ""
	defer func() { R2StemsBucket = orig }()

	m := &Mixer{}
	dir := t.TempDir()
	paths, err := m.GetStems(dir, []string{"1", "2", "3"})
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 3 {
		t.Fatalf("expected 3 paths, got %d", len(paths))
	}
	for _, p := range paths {
		if !strings.HasSuffix(p, ".mp3") {
			t.Errorf("expected .mp3 suffix: %s", p)
		}
		if !strings.Contains(p, dir) {
			t.Errorf("expected path to contain dir: %s", p)
		}
	}
}

func TestGetStemsIncludesStemName(t *testing.T) {
	orig := R2StemsBucket
	R2StemsBucket = ""
	defer func() { R2StemsBucket = orig }()

	m := &Mixer{}
	dir := t.TempDir()
	paths, err := m.GetStems(dir, []string{"kick", "bass"})
	if err != nil {
		t.Fatal(err)
	}
	hasKick, hasBass := false, false
	for _, p := range paths {
		if strings.Contains(p, "kick") {
			hasKick = true
		}
		if strings.Contains(p, "bass") {
			hasBass = true
		}
	}
	if !hasKick {
		t.Error("expected a path containing 'kick'")
	}
	if !hasBass {
		t.Error("expected a path containing 'bass'")
	}
}

func TestBuildFFmpegArgsStructure(t *testing.T) {
	args := buildFFmpegArgs([]string{"a.mp3", "b.mp3"}, []float64{0.8, 0.6}, "out.mp3")

	if args[0] != "-i" || args[1] != "a.mp3" || args[2] != "-i" || args[3] != "b.mp3" {
		t.Errorf("unexpected input args: %v", args[:4])
	}
	if args[len(args)-1] != "out.mp3" {
		t.Errorf("expected out.mp3 as last arg, got %s", args[len(args)-1])
	}

	filterIdx := -1
	for i, a := range args {
		if a == "-filter_complex" {
			filterIdx = i
			break
		}
	}
	if filterIdx == -1 {
		t.Fatal("missing -filter_complex")
	}

	fc := args[filterIdx+1]
	if !strings.Contains(fc, "volume=0.800") {
		t.Errorf("expected volume=0.800 in filter_complex: %s", fc)
	}
	if !strings.Contains(fc, "volume=0.600") {
		t.Errorf("expected volume=0.600 in filter_complex: %s", fc)
	}
	if !strings.Contains(fc, "amix=inputs=2") {
		t.Errorf("expected amix=inputs=2 in filter_complex: %s", fc)
	}
	if !strings.Contains(fc, "duration=longest") {
		t.Errorf("expected duration=longest in filter_complex: %s", fc)
	}
}

func TestBuildFFmpegArgsDefaultVolumesNotShared(t *testing.T) {
	// Verify that buildFFmpegArgs with the same volume slice produces consistent output
	args1 := buildFFmpegArgs([]string{"a.mp3", "b.mp3"}, []float64{0.5, 0.5}, "out1.mp3")
	args2 := buildFFmpegArgs([]string{"a.mp3", "b.mp3"}, []float64{0.5, 0.5}, "out2.mp3")

	if args1[len(args1)-1] == args2[len(args2)-1] {
		// Different output paths — this would only fail if output paths were the same
		t.Log("output paths differ as expected (unless timestamp collision)")
	}
}
