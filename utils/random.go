package utils

import (
	"math/rand"
)

func RandomInt32Range(min, max int32) int32 {
	if max <= min {
		return min
	}
	return min + rand.Int31n(max-min+1)
}
