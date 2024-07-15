package utils

import (
	"fmt"
	"crypto/rand"
    "time"
)

type ResetCode struct {
	Email     string
	Code      string
	ExpiresAt time.Time
}

var ResetCodes = make(map[string]ResetCode)

func GenerateResetCode() string {
    code := make([]byte, 3)
    rand.Read(code)
    return fmt.Sprintf("%06d", int(code[0])<<16|int(code[1])<<8|int(code[2]))
}

