package sb

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestParallelSummary_AllPassed(t *testing.T) {
	results := []ParallelResult{
		{Name: "test1", Passed: true, Duration: 100 * time.Millisecond},
		{Name: "test2", Passed: true, Duration: 200 * time.Millisecond},
	}
	summary := ParallelSummary(results)
	if !strings.Contains(summary, "2 passed, 0 failed, 2 total") {
		t.Errorf("unexpected summary: %s", summary)
	}
}

func TestParallelSummary_Mixed(t *testing.T) {
	results := []ParallelResult{
		{Name: "test1", Passed: true, Duration: 100 * time.Millisecond},
		{Name: "test2", Passed: false, Duration: 50 * time.Millisecond, Error: errors.New("fail")},
		{Name: "test3", Passed: true, Duration: 150 * time.Millisecond},
	}
	summary := ParallelSummary(results)
	if !strings.Contains(summary, "2 passed, 1 failed, 3 total") {
		t.Errorf("unexpected summary: %s", summary)
	}
}

func TestParallelSummary_Empty(t *testing.T) {
	summary := ParallelSummary(nil)
	if !strings.Contains(summary, "0 passed, 0 failed, 0 total") {
		t.Errorf("unexpected summary: %s", summary)
	}
}

func TestParallelResult_Fields(t *testing.T) {
	err := errors.New("something went wrong")
	r := ParallelResult{
		Name:     "my_test",
		Passed:   false,
		Duration: 500 * time.Millisecond,
		Error:    err,
	}
	if r.Name != "my_test" {
		t.Errorf("expected name my_test, got %s", r.Name)
	}
	if r.Passed {
		t.Error("expected Passed=false")
	}
	if r.Error != err {
		t.Error("expected stored error")
	}
}

func TestTestFunc_Fields(t *testing.T) {
	tf := TestFunc{
		Name: "login_test",
		Fn: func(p *Page) error {
			return nil
		},
	}
	if tf.Name != "login_test" {
		t.Errorf("expected name login_test, got %s", tf.Name)
	}
	if tf.Fn == nil {
		t.Error("expected non-nil function")
	}
}
