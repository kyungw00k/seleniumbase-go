package sb

import (
	"encoding/xml"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func sampleResults() []ParallelResult {
	return []ParallelResult{
		{Name: "TestLogin", Passed: true, Duration: 2500 * time.Millisecond},
		{Name: "TestCheckout", Passed: false, Duration: 1200 * time.Millisecond, Error: errors.New("element not found: #cart")},
		{Name: "TestSearch", Passed: true, Duration: 800 * time.Millisecond},
	}
}

func TestGenerateJUnitReport(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.xml")

	err := GenerateJUnitReport(path, "MySuite", sampleResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read report: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "<?xml") {
		t.Error("expected XML header")
	}
	if !strings.Contains(content, `name="MySuite"`) {
		t.Error("expected suite name")
	}
	if !strings.Contains(content, `tests="3"`) {
		t.Error("expected tests=3")
	}
	if !strings.Contains(content, `failures="1"`) {
		t.Error("expected failures=1")
	}
	if !strings.Contains(content, `name="TestLogin"`) {
		t.Error("expected TestLogin test case")
	}
	if !strings.Contains(content, "element not found") {
		t.Error("expected failure message")
	}

	// Verify it's valid XML
	var suites JUnitTestSuites
	if err := xml.Unmarshal(data, &suites); err != nil {
		t.Fatalf("invalid XML: %v", err)
	}
	if len(suites.Suites) != 1 {
		t.Fatalf("expected 1 suite, got %d", len(suites.Suites))
	}
	suite := suites.Suites[0]
	if suite.Tests != 3 {
		t.Errorf("expected 3 tests, got %d", suite.Tests)
	}
	if suite.Failures != 1 {
		t.Errorf("expected 1 failure, got %d", suite.Failures)
	}
}

func TestGenerateJUnitReport_Empty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.xml")

	err := GenerateJUnitReport(path, "Empty", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read report: %v", err)
	}

	var suites JUnitTestSuites
	if err := xml.Unmarshal(data, &suites); err != nil {
		t.Fatalf("invalid XML: %v", err)
	}
	if suites.Suites[0].Tests != 0 {
		t.Errorf("expected 0 tests, got %d", suites.Suites[0].Tests)
	}
}

func TestGenerateHTMLReport(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.html")

	err := GenerateHTMLReport(path, "Test Report", sampleResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read report: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "<!DOCTYPE html>") {
		t.Error("expected HTML doctype")
	}
	if !strings.Contains(content, "Test Report") {
		t.Error("expected report title")
	}
	if !strings.Contains(content, "TestLogin") {
		t.Error("expected TestLogin in report")
	}
	if !strings.Contains(content, "PASS") {
		t.Error("expected PASS status")
	}
	if !strings.Contains(content, "FAIL") {
		t.Error("expected FAIL status")
	}
	if !strings.Contains(content, "element not found") {
		t.Error("expected error message in report")
	}
	if !strings.Contains(content, "seleniumbase-go") {
		t.Error("expected footer")
	}
}

func TestGenerateHTMLReport_Empty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.html")

	err := GenerateHTMLReport(path, "Empty Report", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read report: %v", err)
	}

	if !strings.Contains(string(data), "Empty Report") {
		t.Error("expected report title")
	}
}

func TestFormatDuration_Milliseconds(t *testing.T) {
	d := 500 * time.Millisecond
	s := formatDuration(d)
	if s != "500ms" {
		t.Errorf("expected 500ms, got %s", s)
	}
}

func TestFormatDuration_Seconds(t *testing.T) {
	d := 2500 * time.Millisecond
	s := formatDuration(d)
	if s != "2.50s" {
		t.Errorf("expected 2.50s, got %s", s)
	}
}

func TestJUnitReport_InvalidPath(t *testing.T) {
	err := GenerateJUnitReport("/nonexistent/dir/report.xml", "test", nil)
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestHTMLReport_InvalidPath(t *testing.T) {
	err := GenerateHTMLReport("/nonexistent/dir/report.html", "test", nil)
	if err == nil {
		t.Error("expected error for invalid path")
	}
}
