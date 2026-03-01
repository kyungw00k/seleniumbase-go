package sb

import (
	"strings"
	"testing"
)

func TestParseActions_Valid(t *testing.T) {
	raw := []interface{}{
		[]interface{}{"click", "#btn", "", float64(1000)},
		[]interface{}{"input", "#name", "Alice", float64(1001)},
	}
	actions, err := parseActions(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(actions) != 2 {
		t.Fatalf("expected 2 actions, got %d", len(actions))
	}
	if actions[0].Action != "click" || actions[0].Selector != "#btn" {
		t.Errorf("action[0] = %+v, want click on #btn", actions[0])
	}
	if actions[1].Action != "input" || actions[1].Data != "Alice" {
		t.Errorf("action[1] = %+v, want input Alice", actions[1])
	}
}

func TestParseActions_Empty(t *testing.T) {
	raw := []interface{}{}
	actions, err := parseActions(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(actions) != 0 {
		t.Errorf("expected 0 actions, got %d", len(actions))
	}
}

func TestParseActions_Invalid(t *testing.T) {
	_, err := parseActions("not an array")
	if err == nil {
		t.Error("expected error for non-array input")
	}
}

func TestGenerateGoCode_Empty(t *testing.T) {
	code := GenerateGoCode(nil)
	if !strings.Contains(code, "return nil") {
		t.Error("empty actions should produce function with return nil")
	}
	if strings.Contains(code, "p.") {
		t.Error("empty actions should not contain any method calls")
	}
}

func TestGenerateGoCode_Open(t *testing.T) {
	actions := []RecordedAction{
		{Action: "begin", Data: "https://example.com"},
	}
	code := GenerateGoCode(actions)
	if !strings.Contains(code, `p.Open("https://example.com")`) {
		t.Errorf("expected Open call, got:\n%s", code)
	}
}

func TestGenerateGoCode_Click(t *testing.T) {
	actions := []RecordedAction{
		{Action: "click", Selector: "#btn"},
	}
	code := GenerateGoCode(actions)
	if !strings.Contains(code, `p.Click("#btn")`) {
		t.Errorf("expected Click call, got:\n%s", code)
	}
}

func TestGenerateGoCode_DoubleClick(t *testing.T) {
	actions := []RecordedAction{
		{Action: "dbclk", Selector: "#item"},
	}
	code := GenerateGoCode(actions)
	if !strings.Contains(code, `p.DoubleClick("#item")`) {
		t.Errorf("expected DoubleClick call, got:\n%s", code)
	}
}

func TestGenerateGoCode_Type(t *testing.T) {
	actions := []RecordedAction{
		{Action: "input", Selector: "#name", Data: "Alice"},
	}
	code := GenerateGoCode(actions)
	if !strings.Contains(code, `p.Type("#name", "Alice")`) {
		t.Errorf("expected Type call, got:\n%s", code)
	}
}

func TestGenerateGoCode_TypeTrimsNewline(t *testing.T) {
	actions := []RecordedAction{
		{Action: "input", Selector: "#name", Data: "Alice\n"},
	}
	code := GenerateGoCode(actions)
	if !strings.Contains(code, `p.Type("#name", "Alice")`) {
		t.Errorf("expected Type call with trimmed newline, got:\n%s", code)
	}
}

func TestGenerateGoCode_SelectOption(t *testing.T) {
	actions := []RecordedAction{
		{Action: "s_opt", Selector: "#color", Data: "blue"},
	}
	code := GenerateGoCode(actions)
	if !strings.Contains(code, `p.SelectOption("#color", "blue")`) {
		t.Errorf("expected SelectOption call, got:\n%s", code)
	}
}

func TestGenerateGoCode_AssertText(t *testing.T) {
	actions := []RecordedAction{
		{Action: "as_te", Selector: "h1", Data: "Welcome"},
	}
	code := GenerateGoCode(actions)
	if !strings.Contains(code, `p.AssertText("Welcome", "h1")`) {
		t.Errorf("expected AssertText call, got:\n%s", code)
	}
}

func TestGenerateGoCode_AssertElement(t *testing.T) {
	actions := []RecordedAction{
		{Action: "as_el", Selector: "#logo"},
	}
	code := GenerateGoCode(actions)
	if !strings.Contains(code, `p.AssertElement("#logo")`) {
		t.Errorf("expected AssertElement call, got:\n%s", code)
	}
}

func TestGenerateGoCode_Highlight(t *testing.T) {
	actions := []RecordedAction{
		{Action: "hi_lt", Selector: ".item"},
	}
	code := GenerateGoCode(actions)
	if !strings.Contains(code, `p.Highlight(".item")`) {
		t.Errorf("expected Highlight call, got:\n%s", code)
	}
}

func TestGenerateGoCode_Screenshot(t *testing.T) {
	actions := []RecordedAction{
		{Action: "savsc", Selector: "", Data: ""},
	}
	code := GenerateGoCode(actions)
	if !strings.Contains(code, `p.Screenshot(`) {
		t.Errorf("expected Screenshot call, got:\n%s", code)
	}
}

func TestGenerateGoCode_ClickThenDoubleClick(t *testing.T) {
	actions := []RecordedAction{
		{Action: "click", Selector: "#item"},
		{Action: "dbclk", Selector: "#item"},
	}
	code := GenerateGoCode(actions)
	if strings.Contains(code, `p.Click("#item")`) {
		t.Error("click followed by dbclk on same selector should not emit Click")
	}
	if !strings.Contains(code, `p.DoubleClick("#item")`) {
		t.Errorf("expected DoubleClick call, got:\n%s", code)
	}
}

func TestGenerateGoCode_ClickThenDoubleClickEmptySelector(t *testing.T) {
	// The JS recorder emits dbclk with an empty selector; codegen should still dedup
	actions := []RecordedAction{
		{Action: "click", Selector: "#item"},
		{Action: "dbclk", Selector: ""},
	}
	code := GenerateGoCode(actions)
	if strings.Contains(code, `p.Click("#item")`) {
		t.Error("click followed by dbclk with empty selector should not emit Click")
	}
	if !strings.Contains(code, `p.DoubleClick("#item")`) {
		t.Errorf("expected DoubleClick with inherited selector, got:\n%s", code)
	}
}

func TestGenerateGoCode_UnknownAction(t *testing.T) {
	actions := []RecordedAction{
		{Action: "xyz123", Selector: "#foo"},
	}
	code := GenerateGoCode(actions)
	if !strings.Contains(code, "// unknown action: xyz123") {
		t.Errorf("expected unknown action comment, got:\n%s", code)
	}
}

func TestEscapeGoString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"no special chars", "hello", "hello"},
		{"double quote", `say "hi"`, `say \"hi\"`},
		{"backslash", `path\to\file`, `path\\to\\file`},
		{"newline", "line1\nline2", `line1\nline2`},
		{"tab", "col1\tcol2", `col1\tcol2`},
		{"mixed", "a\"b\\c\nd", `a\"b\\c\nd`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := escapeGoString(tt.input)
			if got != tt.want {
				t.Errorf("escapeGoString(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
