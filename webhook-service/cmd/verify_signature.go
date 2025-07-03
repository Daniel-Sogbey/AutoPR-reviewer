package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"strings"
)

func computeHMAC(reqBody []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(reqBody)
	sha256Hash := h.Sum(nil)
	return hex.EncodeToString(sha256Hash)
}

func VerifySignature(reqBody []byte, secret, signatureHeader string) bool {
	if !strings.HasPrefix(signatureHeader, "sha256=") {
		return false
	}

	sentSignature := strings.TrimPrefix(signatureHeader, "sha256=")

	expectedSignature := computeHMAC(reqBody, secret)

	return subtle.ConstantTimeCompare([]byte(sentSignature), []byte(expectedSignature)) == 1
}
