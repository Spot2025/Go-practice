package reviewpathdemo

import (
	"errors"
	"testing"
)

func TestValidateOrderRejectsIncompleteInput(t *testing.T) {
	tests := []struct {
		name  string
		order Order
	}{
		{name: "missing order ID", order: validOrderWith(func(order *Order) { order.ID = "" })},
		{name: "missing customer", order: validOrderWith(func(order *Order) { order.CustomerID = "" })},
		{name: "invalid country", order: validOrderWith(func(order *Order) { order.DestinationCountry = "USA" })},
		{name: "missing items", order: validOrderWith(func(order *Order) { order.Items = nil })},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := validateOrder(test.order); !errors.Is(err, ErrInvalidOrder) {
				t.Fatalf("validateOrder() error = %v, want ErrInvalidOrder", err)
			}
		})
	}
}

func validOrderWith(change func(*Order)) Order {
	order := Order{
		ID:                 "ORDER-1",
		CustomerID:         "CUSTOMER-1",
		DestinationCountry: "US",
		Items: []Item{
			{SKU: "BOOK", PriceCents: 1000, Quantity: 1, WeightGrams: 100},
		},
	}
	change(&order)
	return order
}
