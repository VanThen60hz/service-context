package core

import (
	"math/rand"
	"time"
)

func GenerateOTP() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	digits := "0123456789"
	otp := ""
	for i := 0; i < 6; i++ {
		otp += string(digits[r.Intn(len(digits))])
	}
	return otp
}
