package ordoe2e

import "testing"

func TestCheckoutLabel(t *testing.T) {
	if got := CheckoutLabel(120, 15); got != "total: 135" {
		t.Fatalf("CheckoutLabel() = %q, want %q", got, "total: 135")
	}
}
