package sb

import "testing"

func TestGetConsoleErrors(t *testing.T) {
	p := &Page{
		consoleMessages: []ConsoleMessage{
			{Type: "log", Text: "hello"},
			{Type: "error", Text: "something broke"},
			{Type: "warning", Text: "deprecated"},
			{Type: "error", Text: "another error"},
		},
	}
	errors := p.GetConsoleErrors()
	if len(errors) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errors))
	}
	if errors[0].Text != "something broke" {
		t.Errorf("expected 'something broke', got %q", errors[0].Text)
	}
	if errors[1].Text != "another error" {
		t.Errorf("expected 'another error', got %q", errors[1].Text)
	}
}

func TestGetConsoleErrors_Empty(t *testing.T) {
	p := &Page{
		consoleMessages: []ConsoleMessage{
			{Type: "log", Text: "hello"},
		},
	}
	errors := p.GetConsoleErrors()
	if len(errors) != 0 {
		t.Fatalf("expected 0 errors, got %d", len(errors))
	}
}

func TestAssertNoJsErrors_Pass(t *testing.T) {
	p := &Page{
		consoleMessages: []ConsoleMessage{
			{Type: "log", Text: "info"},
			{Type: "warning", Text: "warn"},
		},
	}
	if err := p.AssertNoJsErrors(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestAssertNoJsErrors_Fail(t *testing.T) {
	p := &Page{
		consoleMessages: []ConsoleMessage{
			{Type: "error", Text: "ReferenceError: x is not defined"},
		},
	}
	err := p.AssertNoJsErrors()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetConsoleMessages(t *testing.T) {
	msgs := []ConsoleMessage{
		{Type: "log", Text: "a"},
		{Type: "error", Text: "b"},
	}
	p := &Page{consoleMessages: msgs}
	got := p.GetConsoleMessages()
	if len(got) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(got))
	}
}
