package utils

import (
	"math/rand"
	"time"
)

func NewRandomString(size int) string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	runes := make([]rune, size)
	for i := range runes {
		runes[i] = chars[random.Intn(len(chars))]
	}

	return string(runes)
}
