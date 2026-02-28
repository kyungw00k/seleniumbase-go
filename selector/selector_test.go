package selector

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// CSS selectors pass through unchanged
		{"id css", "#id", "#id"},
		{"class css", ".class", ".class"},
		{"descendant css", "div > span", "div > span"},
		{"attribute css", "[attr=val]", "[attr=val]"},
		{"tag css", "div", "div"},
		{"complex css", "div.container > ul li.active", "div.container > ul li.active"},

		// XPath pass through
		{"absolute xpath", "//div", "//div"},
		{"relative xpath", "./div", "./div"},
		{"grouped xpath", "(//div)[1]", "(//div)[1]"},
		{"xpath prefix", "xpath=//div", "//div"},
		{"xpath prefix with predicate", "xpath=//input[@name='email']", "//input[@name='email']"},

		// link= conversion
		{"link text", "link=Login", `a:has-text("Login")`},
		{"link text with spaces", "link=Sign In", `a:has-text("Sign In")`},
		{"link text special chars", "link=Hello & World", `a:has-text("Hello & World")`},

		// partial_link= conversion
		{"partial link text", "partial_link=Log", `a:has-text("Log")`},
		{"partial link text with spaces", "partial_link=Sign", `a:has-text("Sign")`},

		// name= conversion
		{"name selector", "name=email", `[name="email"]`},
		{"name selector with underscore", "name=first_name", `[name="first_name"]`},
		{"name selector with hyphen", "name=user-id", `[name="user-id"]`},

		// css= prefix stripping
		{"css prefix", "css=div.x", "div.x"},
		{"css prefix complex", "css=#main .content", "#main .content"},
		{"css prefix attribute", "css=[data-testid=button]", "[data-testid=button]"},

		// id= conversion
		{"id prefix", "id=main", "#main"},
		{"id prefix with numbers", "id=item-123", "#item-123"},

		// text= pass through (Playwright supports this)
		{"text selector", "text=Hello", "text=Hello"},
		{"text selector with spaces", "text=Hello World", "text=Hello World"},

		// Edge cases
		{"empty string", "", ""},
		{"whitespace", "   ", "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.input)
			if got != tt.want {
				t.Errorf("Parse(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsXPath(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"absolute xpath", "//div", true},
		{"relative xpath", "./span", true},
		{"grouped xpath", "(//li)[1]", true},
		{"xpath= prefix", "xpath=//div", true},
		{"xpath= with predicate", "xpath=//input[@id='email']", true},
		{"css selector", "#id", false},
		{"class selector", ".class", false},
		{"tag selector", "div", false},
		{"link selector", "link=Login", false},
		{"name selector", "name=email", false},
		{"empty string", "", false},
		{"text selector", "text=Hello", false},
		{"id prefix", "id=main", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsXPath(tt.input)
			if got != tt.want {
				t.Errorf("IsXPath(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsLinkText(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"link= prefix", "link=Login", true},
		{"link= with spaces", "link=Sign In", true},
		{"partial_link= prefix", "partial_link=Log", true},
		{"partial_link= with spaces", "partial_link=Sign", true},
		{"css selector", "#id", false},
		{"xpath", "//div", false},
		{"name selector", "name=email", false},
		{"text selector", "text=Login", false},
		{"id prefix", "id=main", false},
		{"empty string", "", false},
		{"link without prefix", "Login", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsLinkText(tt.input)
			if got != tt.want {
				t.Errorf("IsLinkText(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
