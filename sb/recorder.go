package sb

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/playwright-community/playwright-go"
)

//go:embed recorder.js
var recorderJS string

// RecordedAction represents a single recorded browser action.
type RecordedAction struct {
	Action    string  // action code (click, input, dbclk, etc.)
	Selector  string  // CSS selector
	Data      string  // associated data (text input, URL, etc.)
	Timestamp float64 // millisecond timestamp
}

// parseActions converts raw sessionStorage JSON into typed actions.
func parseActions(raw interface{}) ([]RecordedAction, error) {
	outer, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("sb: expected array of actions, got %T", raw)
	}
	actions := make([]RecordedAction, 0, len(outer))
	for _, item := range outer {
		arr, ok := item.([]interface{})
		if !ok || len(arr) < 4 {
			continue
		}
		action, _ := arr[0].(string)
		selector, _ := arr[1].(string)
		data, _ := arr[2].(string)
		ts, _ := arr[3].(float64)
		actions = append(actions, RecordedAction{
			Action:    action,
			Selector:  selector,
			Data:      data,
			Timestamp: ts,
		})
	}
	return actions, nil
}

// StartRecording injects the recorder JavaScript into the page and all future navigations.
func (p *Page) StartRecording() error {
	if p.recording {
		return nil
	}
	if err := p.context.AddInitScript(playwright.Script{Content: &recorderJS}); err != nil {
		return fmt.Errorf("sb: could not inject recorder: %w", err)
	}
	if _, err := p.pw.Evaluate(recorderJS); err != nil {
		return fmt.Errorf("sb: could not start recorder on current page: %w", err)
	}
	p.recording = true
	return nil
}

// StopRecording stops the recorder and returns all captured actions.
func (p *Page) StopRecording() ([]RecordedAction, error) {
	if !p.recording {
		return nil, nil
	}
	raw, err := p.pw.Evaluate(`JSON.parse(sessionStorage.getItem('recorded_actions') || '[]')`)
	if err != nil {
		return nil, fmt.Errorf("sb: could not read recorded actions: %w", err)
	}
	p.recording = false
	return parseActions(raw)
}

// GetRecordedActions reads all captured actions without stopping the recorder.
func (p *Page) GetRecordedActions() ([]RecordedAction, error) {
	raw, err := p.pw.Evaluate(`JSON.parse(sessionStorage.getItem('recorded_actions') || '[]')`)
	if err != nil {
		return nil, fmt.Errorf("sb: could not read recorded actions: %w", err)
	}
	return parseActions(raw)
}

// RunRecorder opens a browser with recording enabled, waits for it to close,
// and writes the generated Go code to outputPath.
func RunRecorder(outputPath string, opts ...Option) error {
	return Run(func(p *Page) error {
		if err := p.StartRecording(); err != nil {
			return err
		}
		p.pw.WaitForEvent("close")
		actions, err := p.StopRecording()
		if err != nil {
			return err
		}
		code := GenerateGoCode(actions)
		return os.WriteFile(outputPath, []byte(code), 0644)
	}, opts...)
}
