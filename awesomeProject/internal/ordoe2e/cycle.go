// Package ordoe2e contains small fixtures used to exercise Ordo in pull requests.
package ordoe2e

// ParityLabel describes whether a non-negative integer is even or odd.
func ParityLabel(value int) string {
	if Even(value) {
		return "even"
	}

	return "odd"
}

// Even reports whether a non-negative integer is even.
func Even(value int) bool {
	if value == 0 {
		return true
	}

	return Odd(value - 1)
}

// Odd reports whether a non-negative integer is odd.
func Odd(value int) bool {
	if value == 0 {
		return false
	}

	return Even(value - 1)
}
