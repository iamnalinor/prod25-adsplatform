package floatutil

import "math"

// Norm normalizes target into [-1;1] scale in scope of only two values.
// It is equivalent to target / max(abs(target), abs(other)) with exception to the case
// where both values are 0.
// Example: Norm(33, 100) => 0.33, Norm(100, 33) => 1.
func Norm(target, other float64) float64 {
	if max(target, other) == 0 {
		return 0
	}
	return target / max(math.Abs(target), math.Abs(other))
}
