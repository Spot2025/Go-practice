package reviewpathdemo

import "context"

type Service struct {
	inventory Inventory
	payments  PaymentGateway
}

func NewService(inventory Inventory, payments PaymentGateway) *Service {
	return &Service{
		inventory: inventory,
		payments:  payments,
	}
}

func (s *Service) Checkout(ctx context.Context, input Order) (Result, error) {
	order := normalizeOrder(input)
	if err := validateOrder(order); err != nil {
		return Result{}, CheckoutError{Stage: "validation", Err: err}
	}

	if err := reserveOrder(ctx, s.inventory, order); err != nil {
		return Result{}, CheckoutError{Stage: "inventory", Err: err}
	}

	totals := calculateTotals(order)
	payment, err := authorizeOrderPayment(ctx, s.payments, order, totals)
	if err != nil {
		return Result{}, CheckoutError{Stage: "payment", Err: err}
	}

	return assembleResult(order, totals, payment), nil
}

func calculateTotals(order Order) Totals {
	subtotal := calculateSubtotal(order.Items)
	discount := calculateDiscount(order, subtotal)
	taxable := taxableAmount(subtotal, discount)
	tax := calculateTax(order.DestinationCountry, taxable)
	shipping := calculateShipping(order)

	return Totals{
		Subtotal: subtotal,
		Discount: discount,
		Tax:      tax,
		Shipping: shipping,
		Grand:    sumMoney(taxable, tax, shipping),
	}
}

func assembleResult(order Order, totals Totals, payment string) Result {
	return Result{
		Receipt: buildReceipt(order, totals, payment),
		Audit:   buildAuditTrail(order, totals, payment),
	}
}
