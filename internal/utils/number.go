package utils

import "math"

func ToFixed(x float64, y int) float64 {
	z := math.Pow(10, float64(y))
	return math.Floor(x*z) / z
}
