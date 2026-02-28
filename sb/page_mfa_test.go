package sb

import "testing"

func TestGenerateTOTPAt(t *testing.T) {
	// RFC 6238 test vector: SHA1, secret "12345678901234567890" (ASCII)
	// base32("12345678901234567890") = "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ"
	secret := "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ"

	tests := []struct {
		name     string
		unixTime int64
		want     string
	}{
		{"T=59", 59, "287082"},
		{"T=1111111109", 1111111109, "081804"},
		{"T=1111111111", 1111111111, "050471"},
		{"T=1234567890", 1234567890, "005924"},
		{"T=2000000000", 2000000000, "279037"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateTOTPAt(secret, tt.unixTime)
			if err != nil {
				t.Fatalf("generateTOTPAt() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("generateTOTPAt(%q, %d) = %q, want %q", secret, tt.unixTime, got, tt.want)
			}
		})
	}
}

func TestGenerateTOTP_InvalidSecret(t *testing.T) {
	_, err := generateTOTP("!!!invalid!!!")
	if err == nil {
		t.Error("expected error for invalid base32 secret, got nil")
	}
}

func TestGenerateTOTP_SpacesAndCase(t *testing.T) {
	// Spaces and lowercase should be normalized to the same key
	secret1 := "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ"
	secret2 := "gezd gnbv gy3t qojq gezd gnbv gy3t qojq"

	code1, err := generateTOTPAt(secret1, 59)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	code2, err := generateTOTPAt(secret2, 59)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Both should produce the same code since they normalize to the same key
	if code1 != code2 {
		t.Errorf("expected same code for equivalent secrets, got %q and %q", code1, code2)
	}
}

func TestGenerateTOTP_OutputFormat(t *testing.T) {
	// Verify the output is always 6 digits with leading zeros
	code, err := generateTOTPAt("GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ", 59)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(code) != 6 {
		t.Errorf("expected 6-digit code, got %q (len %d)", code, len(code))
	}
	for _, c := range code {
		if c < '0' || c > '9' {
			t.Errorf("expected only digits, got %q", code)
			break
		}
	}
}
