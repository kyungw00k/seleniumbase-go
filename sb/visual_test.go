package sb

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func createTestPNG(t *testing.T, path string, width, height int, c color.Color) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("could not create test PNG: %v", err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("could not encode test PNG: %v", err)
	}
}

func TestCompareImages_Identical(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "base.png")
	latest := filepath.Join(dir, "latest.png")
	diff := filepath.Join(dir, "diff.png")

	createTestPNG(t, base, 100, 100, color.RGBA{255, 0, 0, 255})
	createTestPNG(t, latest, 100, 100, color.RGBA{255, 0, 0, 255})

	result, err := compareImages(base, latest, diff)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.DiffPixels != 0 {
		t.Errorf("expected 0 diff pixels, got %d", result.DiffPixels)
	}
	if result.DiffPercent != 0 {
		t.Errorf("expected 0%% diff, got %.2f%%", result.DiffPercent)
	}
	if result.TotalPixels != 10000 {
		t.Errorf("expected 10000 total pixels, got %d", result.TotalPixels)
	}
}

func TestCompareImages_Different(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "base.png")
	latest := filepath.Join(dir, "latest.png")
	diff := filepath.Join(dir, "diff.png")

	createTestPNG(t, base, 100, 100, color.RGBA{255, 0, 0, 255})
	createTestPNG(t, latest, 100, 100, color.RGBA{0, 0, 255, 255})

	result, err := compareImages(base, latest, diff)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.DiffPixels != 10000 {
		t.Errorf("expected 10000 diff pixels, got %d", result.DiffPixels)
	}
	if result.DiffPercent != 100 {
		t.Errorf("expected 100%% diff, got %.2f%%", result.DiffPercent)
	}
	// Diff image should be created
	if _, err := os.Stat(diff); os.IsNotExist(err) {
		t.Error("diff image was not created")
	}
}

func TestCompareImages_DifferentSizes(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "base.png")
	latest := filepath.Join(dir, "latest.png")
	diff := filepath.Join(dir, "diff.png")

	createTestPNG(t, base, 100, 100, color.RGBA{255, 0, 0, 255})
	createTestPNG(t, latest, 200, 150, color.RGBA{255, 0, 0, 255})

	result, err := compareImages(base, latest, diff)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Extra pixels should count as differences
	if result.DiffPixels == 0 {
		t.Error("expected diff pixels for different-sized images")
	}
	if result.TotalPixels != 200*150 {
		t.Errorf("expected %d total pixels, got %d", 200*150, result.TotalPixels)
	}
}

func TestColorDistance_Same(t *testing.T) {
	c := color.RGBA{128, 64, 32, 255}
	d := colorDistance(c, c)
	if d != 0 {
		t.Errorf("expected 0 distance for same color, got %f", d)
	}
}

func TestColorDistance_Different(t *testing.T) {
	c1 := color.RGBA{0, 0, 0, 255}
	c2 := color.RGBA{255, 255, 255, 255}
	d := colorDistance(c1, c2)
	if d <= 0 {
		t.Error("expected positive distance for different colors")
	}
}

func TestVisualResult_MatchThreshold(t *testing.T) {
	r := &VisualResult{DiffPercent: 0.5, DiffPixels: 50, TotalPixels: 10000}
	r.Match = r.DiffPercent <= 1.0
	if !r.Match {
		t.Error("0.5% should match with 1.0% threshold")
	}

	r.Match = r.DiffPercent <= 0.1
	if r.Match {
		t.Error("0.5% should not match with 0.1% threshold")
	}
}

func TestLoadPNG_NotFound(t *testing.T) {
	_, err := loadPNG("/nonexistent/path.png")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}
