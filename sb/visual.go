package sb

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"

	"github.com/playwright-community/playwright-go"
)

// VisualResult holds the result of a visual comparison.
type VisualResult struct {
	Match       bool    // true if images match within threshold
	DiffPixels  int     // number of pixels that differ
	TotalPixels int     // total pixel count
	DiffPercent float64 // percentage of pixels that differ
	DiffPath    string  // path to the diff image (if generated)
}

// defaultBaselineDir returns the default directory for visual baselines.
func defaultBaselineDir() string {
	return "visual_baseline"
}

// CheckWindow captures a screenshot and compares it against a stored baseline.
// On first run (no baseline exists), it saves the screenshot as the new baseline.
// Returns a VisualResult with comparison details.
//
// name: identifier for this visual check (used as filename)
// threshold: maximum allowed percentage of differing pixels (0.0 = exact match, 1.0 = up to 1%)
func (p *Page) CheckWindow(name string, threshold float64) (*VisualResult, error) {
	baseDir := defaultBaselineDir()
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("sb: could not create baseline dir: %w", err)
	}

	baselinePath := filepath.Join(baseDir, name+"_baseline.png")
	latestPath := filepath.Join(baseDir, name+"_latest.png")
	diffPath := filepath.Join(baseDir, name+"_diff.png")

	// Take screenshot of current page
	_, err := p.pw.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String(latestPath),
		FullPage: playwright.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("sb: could not take screenshot: %w", err)
	}

	// If no baseline exists, save current as baseline
	if _, err := os.Stat(baselinePath); os.IsNotExist(err) {
		latestData, err := os.ReadFile(latestPath)
		if err != nil {
			return nil, fmt.Errorf("sb: could not read latest screenshot: %w", err)
		}
		if err := os.WriteFile(baselinePath, latestData, 0644); err != nil {
			return nil, fmt.Errorf("sb: could not save baseline: %w", err)
		}
		return &VisualResult{Match: true, DiffPercent: 0}, nil
	}

	// Compare latest against baseline
	result, err := compareImages(baselinePath, latestPath, diffPath)
	if err != nil {
		return nil, err
	}
	result.Match = result.DiffPercent <= threshold
	return result, nil
}

// UpdateBaseline saves the current page screenshot as the new baseline.
func (p *Page) UpdateBaseline(name string) error {
	baseDir := defaultBaselineDir()
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return fmt.Errorf("sb: could not create baseline dir: %w", err)
	}

	baselinePath := filepath.Join(baseDir, name+"_baseline.png")
	_, err := p.pw.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String(baselinePath),
		FullPage: playwright.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("sb: could not save baseline: %w", err)
	}
	return nil
}

// AssertVisualMatch is an assertion that fails if the visual comparison
// exceeds the given threshold percentage.
func (p *Page) AssertVisualMatch(name string, threshold float64) error {
	result, err := p.CheckWindow(name, threshold)
	if err != nil {
		return err
	}
	if !result.Match {
		return fmt.Errorf("sb: visual mismatch for %q: %.2f%% pixels differ (threshold: %.2f%%)", name, result.DiffPercent, threshold)
	}
	return nil
}

// compareImages compares two PNG images pixel by pixel.
// Returns comparison stats and generates a diff image highlighting differences.
func compareImages(baselinePath, latestPath, diffPath string) (*VisualResult, error) {
	baseImg, err := loadPNG(baselinePath)
	if err != nil {
		return nil, fmt.Errorf("sb: could not load baseline image: %w", err)
	}
	latestImg, err := loadPNG(latestPath)
	if err != nil {
		return nil, fmt.Errorf("sb: could not load latest image: %w", err)
	}

	baseBounds := baseImg.Bounds()
	latestBounds := latestImg.Bounds()

	// Use the larger dimensions for comparison
	maxW := baseBounds.Dx()
	if latestBounds.Dx() > maxW {
		maxW = latestBounds.Dx()
	}
	maxH := baseBounds.Dy()
	if latestBounds.Dy() > maxH {
		maxH = latestBounds.Dy()
	}

	diffImg := image.NewRGBA(image.Rect(0, 0, maxW, maxH))
	totalPixels := maxW * maxH
	diffPixels := 0

	for y := 0; y < maxH; y++ {
		for x := 0; x < maxW; x++ {
			var baseColor, latestColor color.Color
			if x < baseBounds.Dx() && y < baseBounds.Dy() {
				baseColor = baseImg.At(x, y)
			} else {
				baseColor = color.Transparent
			}
			if x < latestBounds.Dx() && y < latestBounds.Dy() {
				latestColor = latestImg.At(x, y)
			} else {
				latestColor = color.Transparent
			}

			if colorDistance(baseColor, latestColor) > 10 {
				diffPixels++
				diffImg.Set(x, y, color.RGBA{255, 0, 0, 200}) // red for differences
			} else {
				// Dim the matching pixels
				r, g, b, a := latestColor.RGBA()
				diffImg.Set(x, y, color.RGBA{
					uint8(r >> 9), uint8(g >> 9), uint8(b >> 9), uint8(a >> 8),
				})
			}
		}
	}

	// Save diff image
	if diffPixels > 0 {
		if err := savePNG(diffPath, diffImg); err != nil {
			return nil, fmt.Errorf("sb: could not save diff image: %w", err)
		}
	}

	diffPercent := 0.0
	if totalPixels > 0 {
		diffPercent = float64(diffPixels) / float64(totalPixels) * 100
	}

	result := &VisualResult{
		DiffPixels:  diffPixels,
		TotalPixels: totalPixels,
		DiffPercent: diffPercent,
	}
	if diffPixels > 0 {
		result.DiffPath = diffPath
	}
	return result, nil
}

// colorDistance calculates the Euclidean distance between two colors.
func colorDistance(c1, c2 color.Color) float64 {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	dr := float64(r1) - float64(r2)
	dg := float64(g1) - float64(g2)
	db := float64(b1) - float64(b2)
	da := float64(a1) - float64(a2)
	return math.Sqrt(dr*dr+dg*dg+db*db+da*da) / 256
}

// loadPNG reads a PNG file and returns the image.
func loadPNG(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}

// savePNG writes an image to a PNG file.
func savePNG(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
