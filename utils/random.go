package utils

import (
	"math/rand"
)

func RandomInt16Range(min, max int16) int16 {
	if max <= min {
		return min
	}
	return min + int16(rand.Int31n(int32(max-min+1)))
}

func RandomInt32Range(min, max int32) int32 {
	if max <= min {
		return min
	}
	return min + rand.Int31n(max-min+1)
}
