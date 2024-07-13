package utils

func Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// Min returns the minimum of two float64 numbers
func Min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
