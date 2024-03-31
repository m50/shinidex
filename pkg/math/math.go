package math

import (
	"math"

	"golang.org/x/exp/constraints"
)

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Ceil(a float64) int {
	i := int(a)
	if a > float64(i) {
		return i+1
	}
	return i
}

func Floor(a float64) int {
	return int(a)
}

func Round(a float64) int {
	return int(math.Round(a))
}