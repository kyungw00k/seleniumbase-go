package sb

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

// generateTOTP generates a 6-digit TOTP code from a base32-encoded secret.
// Implements RFC 6238 with SHA1, 30-second time step, 6-digit output.
func generateTOTP(secret string) (string, error) {
	// Normalize: uppercase and strip spaces
	secret = strings.ToUpper(strings.ReplaceAll(secret, " ", ""))
	// Pad with = if needed
	if m := len(secret) % 8; m != 0 {
		secret += strings.Repeat("=", 8-m)
	}

	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("invalid TOTP secret: %w", err)
	}

	// Time step counter (30-second intervals)
	counter := uint64(time.Now().Unix() / 30)

	// Encode counter as big-endian uint64
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, counter)

	// HMAC-SHA1
	mac := hmac.New(sha1.New, key)
	mac.Write(buf)
	hash := mac.Sum(nil)

	// Dynamic truncation
	offset := hash[len(hash)-1] & 0x0f
	code := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7fffffff

	// 6-digit code
	return fmt.Sprintf("%06d", code%1000000), nil
}

// generateTOTPAt generates a TOTP code for a specific Unix timestamp.
// This is used for testing with known time values.
func generateTOTPAt(secret string, unixTime int64) (string, error) {
	secret = strings.ToUpper(strings.ReplaceAll(secret, " ", ""))
	if m := len(secret) % 8; m != 0 {
		secret += strings.Repeat("=", 8-m)
	}

	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("invalid TOTP secret: %w", err)
	}

	counter := uint64(unixTime / 30)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, counter)

	mac := hmac.New(sha1.New, key)
	mac.Write(buf)
	hash := mac.Sum(nil)

	offset := hash[len(hash)-1] & 0x0f
	code := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7fffffff

	return fmt.Sprintf("%06d", code%1000000), nil
}

// GetMFACode generates a 6-digit TOTP code from the given base32-encoded secret.
// This is a pure computation with no browser interaction.
func (p *Page) GetMFACode(totpKey string) (string, error) {
	return generateTOTP(totpKey)
}

// EnterMFACode generates a TOTP code and types it into the element matching sel,
// then presses Enter.
func (p *Page) EnterMFACode(sel, totpKey string) error {
	code, err := generateTOTP(totpKey)
	if err != nil {
		return err
	}
	if err := p.Type(sel, code); err != nil {
		return err
	}
	return p.Press(sel, "Enter")
}
