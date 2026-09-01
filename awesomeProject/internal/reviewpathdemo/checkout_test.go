package reviewpathdemo

import (
	"context"
	"testing"
)

func TestCheckoutBuildsDependencyAwareResult(t *testing.T) {
	inventory := NewMemoryInventory(map[string]int{
		"COFFEE": 3,
		"MUG":    2,
	})
	service := NewService(inventory, StaticPaymentGateway{Reference: "pay-42"})

	result, err := service.Checkout(context.Background(), Order{
		ID:                 " order-7 ",
		CustomerID:         " customer-9 ",
		DestinationCountry: "us",
		Coupon:             "review10",
		Items: []Item{
			{SKU: " coffee ", PriceCents: 1000, Quantity: 2, WeightGrams: 200},
			{SKU: "mug", PriceCents: 1500, Quantity: 1, WeightGrams: 350},
		},
	})
	if err != nil {
		t.Fatalf("Checkout() error = %v", err)
	}

	if got := result.Receipt.Totals.Grand; got != 3902 {
		t.Fatalf("grand total = %d, want 3902", got)
	}
	if got := result.Receipt.PaymentReference; got != "ORDER-7:pay-42" {
		t.Fatalf("payment reference = %q, want %q", got, "ORDER-7:pay-42")
	}
	if got := len(result.Receipt.Lines); got != 2 {
		t.Fatalf("receipt lines = %d, want 2", got)
	}
	if got := len(result.Audit); got != 5 {
		t.Fatalf("audit entries = %d, want 5", got)
	}
}
