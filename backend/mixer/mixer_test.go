package mixer

import (
	"strings"
	"testing"
)

func TestGetStemsReturnsPaths(t *testing.T) {
	base := t.TempDir()
	m := New(Config{StemsDir: base})

	paths, err := m.GetStems("first-principles", []string{"1", "2", "3"})
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
		if !strings.Contains(p, base) {
			t.Errorf("expected path to contain base dir: %s", p)
		}
		if !strings.Contains(p, "first-principles") {
			t.Errorf("expected path to contain track name: %s", p)
		}
	}
}

func TestGetStemsIncludesStemName(t *testing.T) {
	base := t.TempDir()
	m := New(Config{StemsDir: base})

	paths, err := m.GetStems("test-track", []string{"kick", "bass"})
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
	for _, want := range []string{"volume=0.800", "volume=0.600", "amix=inputs=2", "duration=longest"} {
		if !strings.Contains(fc, want) {
			t.Errorf("expected %q in filter_complex: %s", want, fc)
		}
	}
}
