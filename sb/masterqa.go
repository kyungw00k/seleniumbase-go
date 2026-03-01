package sb

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// MasterQAResult represents the result of a single manual verification step.
type MasterQAResult struct {
	Question   string
	Passed     bool
	Screenshot string // path to screenshot, empty if none
	Timestamp  time.Time
}

// MasterQAPage wraps a Page with manual verification capabilities.
type MasterQAPage struct {
	*Page
	results []MasterQAResult
	reader  *bufio.Reader
}

// Verify pauses execution and asks the tester a yes/no question via stderr/stdin.
// Returns true if the tester answers "y" or "yes".
// In headless mode, automatically returns true with a warning.
func (m *MasterQAPage) Verify(question string) bool {
	if m.config.Headless {
		fmt.Fprintf(os.Stderr, "\n[MasterQA] SKIP (headless): %s\n", question)
		m.results = append(m.results, MasterQAResult{
			Question:  question,
			Passed:    true,
			Timestamp: time.Now(),
		})
		return true
	}

	fmt.Fprintf(os.Stderr, "\n[MasterQA] %s (y/n): ", question)
	answer, _ := m.reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	passed := answer == "y" || answer == "yes"

	result := MasterQAResult{
		Question:  question,
		Passed:    passed,
		Timestamp: time.Now(),
	}

	// Take screenshot for the verification record
	screenshotPath := fmt.Sprintf("masterqa_%d.png", len(m.results)+1)
	if m.pw != nil {
		if err := m.Screenshot(screenshotPath); err == nil {
			result.Screenshot = screenshotPath
		}
	}

	m.results = append(m.results, result)
	return passed
}

// GetResults returns all MasterQA verification results.
func (m *MasterQAPage) GetResults() []MasterQAResult {
	return m.results
}

// SaveReport generates an HTML report of all MasterQA verification results.
func (m *MasterQAPage) SaveReport(path string) error {
	// Convert MasterQAResults to ParallelResults for reuse of existing report generator
	results := make([]ParallelResult, len(m.results))
	for i, r := range m.results {
		pr := ParallelResult{
			Name:     r.Question,
			Passed:   r.Passed,
			Duration: 0,
		}
		if !r.Passed {
			pr.Error = fmt.Errorf("manual verification failed")
		}
		results[i] = pr
	}
	return GenerateHTMLReport(path, "MasterQA Report", results)
}

// RunMasterQA creates a browser session with manual verification support.
func RunMasterQA(t *testing.T, fn func(m *MasterQAPage), opts ...Option) {
	t.Helper()
	err := Run(func(p *Page) error {
		m := &MasterQAPage{
			Page:   p,
			reader: bufio.NewReader(os.Stdin),
		}
		fn(m)

		// Check if any verification failed
		for _, r := range m.results {
			if !r.Passed {
				t.Errorf("[MasterQA] FAILED: %s", r.Question)
			}
		}

		return nil
	}, opts...)
	if err != nil {
		t.Fatal(err)
	}
}
