package handlers

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newID(size int) string {
	alphanumeric := "1234567890abcdefghijklmnopqrstuvwxz"
	var sb strings.Builder
	anlength := len(alphanumeric)
	for i := 0; i < size; i++ {
		c := alphanumeric[rand.Intn(anlength)]
		sb.WriteByte(c)
	}
	return sb.String()
}
