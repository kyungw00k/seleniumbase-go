package sb

import (
	"testing"
	"time"
)

func TestMasterQAResult_Fields(t *testing.T) {
	r := MasterQAResult{
		Question:   "Does the page look correct?",
		Passed:     true,
		Screenshot: "test.png",
		Timestamp:  time.Now(),
	}
	if r.Question != "Does the page look correct?" {
		t.Errorf("expected question, got %s", r.Question)
	}
	if !r.Passed {
		t.Error("expected passed to be true")
	}
	if r.Screenshot != "test.png" {
		t.Errorf("expected test.png, got %s", r.Screenshot)
	}
}

func TestMasterQAPage_GetResults_Empty(t *testing.T) {
	m := &MasterQAPage{
		Page: &Page{config: newDefaultConfig()},
	}
	results := m.GetResults()
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestMasterQAPage_Verify_Headless(t *testing.T) {
	cfg := newDefaultConfig()
	cfg.Headless = true
	m := &MasterQAPage{
		Page: &Page{config: cfg},
	}
	// In headless mode, Verify should return true automatically
	result := m.Verify("Does this look right?")
	if !result {
		t.Error("expected headless Verify to return true")
	}
	if len(m.results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(m.results))
	}
	if !m.results[0].Passed {
		t.Error("expected result to be passed")
	}
}

func TestMasterQAPage_SaveReport(t *testing.T) {
	m := &MasterQAPage{
		Page: &Page{config: newDefaultConfig()},
		results: []MasterQAResult{
			{Question: "Q1", Passed: true, Timestamp: time.Now()},
			{Question: "Q2", Passed: false, Timestamp: time.Now()},
		},
	}
	path := t.TempDir() + "/masterqa_report.html"
	err := m.SaveReport(path)
	if err != nil {
		t.Fatalf("SaveReport failed: %v", err)
	}
}
