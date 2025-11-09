package utils

import (
	"crypto/rand"
	"fmt"
)

// GenerateOTP generates a secure 6-digit numeric OTP as a string.
func GenerateOTP() string {
	// create a 6-digit random number between 000000–999999
	n := make([]byte, 3) // 3 bytes = 24 bits → plenty for 6 digits
	if _, err := rand.Read(n); err != nil {
		// fallback (should rarely happen)
		return "000000"
	}

	// convert bytes to an integer
	num := int(n[0])<<16 | int(n[1])<<8 | int(n[2])
	otp := num % 1000000 // ensure it's 6 digits

	return fmt.Sprintf("%06d", otp)
}
