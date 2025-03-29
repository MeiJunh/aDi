package service

import (
	"aDi/log"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	a, _, err := GenerateToken(1234354)
	log.Debug(a, err)
	//
}

func TestValidateToken(t *testing.T) {
	uid, err := ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDA1ODM1NTUsInVzZXJfaWQiOjEyMzQzNTR9.uVbvwZjJjdcO8oIbXYWwgmljlLsaSFfEiBk6WvlaGZI")
	log.Debug(uid, err)
}
