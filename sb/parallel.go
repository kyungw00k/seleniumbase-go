package sb

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

// TestFunc represents a named test function for parallel execution.
type TestFunc struct {
	Name string
	Fn   func(p *Page) error
}

// ParallelResult holds the result of a single parallel test execution.
type ParallelResult struct {
	Name     string
	Passed   bool
	Duration time.Duration
	Error    error
}

// RunParallel executes multiple test functions concurrently, each with its own
// browser instance. Returns results for all tests.
func RunParallel(tests []TestFunc, opts ...Option) []ParallelResult {
	results := make([]ParallelResult, len(tests))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, runtime.NumCPU())

	for i, tc := range tests {
		wg.Add(1)
		go func(idx int, test TestFunc) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			start := time.Now()
			err := Run(func(p *Page) error {
				return test.Fn(p)
			}, opts...)
			results[idx] = ParallelResult{
				Name:     test.Name,
				Passed:   err == nil,
				Duration: time.Since(start),
				Error:    err,
			}
		}(i, tc)
	}

	wg.Wait()
	return results
}

// RunParallelTest executes multiple test functions concurrently, integrated with testing.T.
// Each test function runs in its own subtest with its own browser instance.
func RunParallelTest(t *testing.T, tests []TestFunc, opts ...Option) {
	t.Helper()
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
			t.Parallel()
			err := Run(func(p *Page) error {
				return tc.Fn(p)
			}, opts...)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

// ParallelSummary returns a formatted summary string from parallel results.
func ParallelSummary(results []ParallelResult) string {
	passed := 0
	failed := 0
	var totalDuration time.Duration
	for _, r := range results {
		if r.Passed {
			passed++
		} else {
			failed++
		}
		totalDuration += r.Duration
	}
	return fmt.Sprintf("%d passed, %d failed, %d total (%.2fs)",
		passed, failed, len(results), totalDuration.Seconds())
}
