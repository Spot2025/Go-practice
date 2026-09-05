// Package ordoe2e contains small fixtures used to exercise Ordo in pull requests.
package ordoe2e

import "fmt"

// CheckoutLabel builds a display label for a checkout total.
func CheckoutLabel(subtotal, shipping int) string {
	total := checkoutTotal(subtotal, shipping)
	return formatTotal(total)
}

func checkoutTotal(subtotal, shipping int) int {
	return normalizeAmount(subtotal) + normalizeAmount(shipping)
}

func normalizeAmount(amount int) int {
	if amount < 0 {
		return 0
	}

	return amount
}

func formatTotal(total int) string {
	return fmt.Sprintf("total: %d", total)
}
