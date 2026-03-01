//go:build integration

package examples

import (
	"os"
	"testing"

	"github.com/kyungw00k/seleniumbase-go/sb"
)

// TestVisualBaseline mirrors xkcd_visual_test.py
// Creates a baseline and then verifies it matches (same page = 0% diff)
func TestVisualBaseline(t *testing.T) {
	// Clean up baselines from previous runs
	os.RemoveAll("visual_baseline")
	defer os.RemoveAll("visual_baseline")

	sb.RunTest(t, func(p *sb.Page) {
		p.Open("https://xkcd.com/554/")

		// First call creates baseline
		result, err := p.CheckWindow("xkcd_554", 0.0)
		if err != nil {
			t.Fatal(err)
		}
		if !result.Match {
			t.Error("first baseline creation should always match")
		}

		// Second call on same page should match exactly
		result, err = p.CheckWindow("xkcd_554", 0.0)
		if err != nil {
			t.Fatal(err)
		}
		if !result.Match {
			t.Errorf("same page should match baseline, got %.2f%% diff", result.DiffPercent)
		}
	}, sb.WithHeadless(true))
}

// TestVisualDOMChange mirrors layout_test.py and test_layout_fail.py
// Creates a baseline, modifies the DOM, then verifies mismatch
func TestVisualDOMChange(t *testing.T) {
	os.RemoveAll("visual_baseline")
	defer os.RemoveAll("visual_baseline")

	sb.RunTest(t, func(p *sb.Page) {
		p.Open("https://xkcd.com/554/")

		// Create baseline
		_, err := p.CheckWindow("xkcd_logo", 0.0)
		if err != nil {
			t.Fatal(err)
		}

		// Modify the logo dimensions (mirrors Python: height 83->130, width 185->120)
		if _, err := p.Evaluate(`document.querySelector('[alt="xkcd.com logo"]').setAttribute("height", "130")`); err != nil {
			t.Fatal(err)
		}
		if _, err := p.Evaluate(`document.querySelector('[alt="xkcd.com logo"]').setAttribute("width", "120")`); err != nil {
			t.Fatal(err)
		}

		// Check with strict threshold — should detect the change
		result, err := p.CheckWindow("xkcd_logo", 0.0)
		if err != nil {
			t.Fatal(err)
		}
		if result.Match {
			t.Error("modified DOM should NOT match baseline with 0% threshold")
		}
		if result.DiffPercent <= 0 {
			t.Error("expected non-zero diff percentage after DOM change")
		}
		t.Logf("Detected %.2f%% pixel difference after DOM change", result.DiffPercent)
	}, sb.WithHeadless(true))
}

// TestVisualAssertMatch mirrors the assertion pattern
func TestVisualAssertMatch(t *testing.T) {
	os.RemoveAll("visual_baseline")
	defer os.RemoveAll("visual_baseline")

	sb.RunTest(t, func(p *sb.Page) {
		p.Open("https://xkcd.com/554/")

		// Create baseline with UpdateBaseline
		if err := p.UpdateBaseline("xkcd_assert"); err != nil {
			t.Fatal(err)
		}

		// AssertVisualMatch should pass on same page
		if err := p.AssertVisualMatch("xkcd_assert", 1.0); err != nil {
			t.Fatalf("AssertVisualMatch should pass on same page: %v", err)
		}
	}, sb.WithHeadless(true))
}
