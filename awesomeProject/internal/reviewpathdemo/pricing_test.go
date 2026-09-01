package reviewpathdemo

import "testing"

func TestCalculateTotals(t *testing.T) {
	tests := []struct {
		name  string
		order Order
		want  Totals
	}{
		{
			name: "discounted domestic order",
			order: Order{
				DestinationCountry: "US",
				Coupon:             "REVIEW20",
				Items: []Item{
					{SKU: "BOOK", PriceCents: 2500, Quantity: 2, WeightGrams: 400},
				},
			},
			want: Totals{Subtotal: 5000, Discount: 1000, Tax: 320, Shipping: 500, Grand: 4820},
		},
		{
			name: "oversized international order",
			order: Order{
				DestinationCountry: "CA",
				Items: []Item{
					{SKU: "CHAIR", PriceCents: 8000, Quantity: 1, WeightGrams: 6000},
				},
			},
			want: Totals{Subtotal: 8000, Shipping: 2100, Grand: 10100},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := calculateTotals(test.order); got != test.want {
				t.Fatalf("calculateTotals() = %+v, want %+v", got, test.want)
			}
		})
	}
}
