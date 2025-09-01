package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixMicro())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}

func RandomString(length int) string {
	var s strings.Builder

	k := len(alphabet)

	for i := 0; i < length; i++ {
		c := alphabet[RandomInt(0, int64(k))]
		s.WriteByte(c)
	}
	return s.String()
}

func RandomOwner() string {
	return RandomString(7)
}

func RandomBalance() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	curr := []string{"EUR", "USD", "INR"}

	return curr[RandomInt(0, int64(len(curr)))]
}
