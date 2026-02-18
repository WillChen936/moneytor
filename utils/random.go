package utils

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomInt16Range(min, max int16) int16 {
	if max <= min {
		return min
	}

	return min + int16(rand.Int31n(int32(max-min+1)))
}

func RandomInt64Range(min, max int64) int64 {
	if max <= min {
		return min
	}

	return min + rand.Int63n(max-min+1)
}

func RandomString(length int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < length; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}
