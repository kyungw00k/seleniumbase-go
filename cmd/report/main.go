package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"github.com/kyungw00k/seleniumbase-go/sb"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: report <junit-xml-files...>")
		os.Exit(1)
	}

	var results []sb.ParallelResult

	for _, path := range os.Args[1:] {
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not read %s: %v\n", path, err)
			continue
		}

		var suites sb.JUnitTestSuites
		if err := xml.Unmarshal(data, &suites); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not parse %s: %v\n", path, err)
			continue
		}

		for _, suite := range suites.Suites {
			for _, tc := range suite.Cases {
				r := sb.ParallelResult{
					Name:     tc.Name,
					Passed:   tc.Failure == nil,
					Duration: time.Duration(tc.Time * float64(time.Second)),
				}
				if tc.Failure != nil {
					r.Error = fmt.Errorf("%s", tc.Failure.Message)
				}
				results = append(results, r)
			}
		}
	}

	if len(results) == 0 {
		fmt.Fprintln(os.Stderr, "no test results found")
		os.Exit(1)
	}

	if err := sb.GenerateHTMLReport("report.html", "SeleniumBase-go Test Report", results); err != nil {
		fmt.Fprintf(os.Stderr, "error generating report: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", sb.ParallelSummary(results))
}
