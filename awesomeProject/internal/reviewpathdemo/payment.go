package reviewpathdemo

import (
	"context"
	"fmt"
)

type PaymentGateway interface {
	Authorize(ctx context.Context, customerID string, amount int64) (string, error)
}

func authorizeOrderPayment(
	ctx context.Context,
	gateway PaymentGateway,
	order Order,
	totals Totals,
) (string, error) {
	reference, err := gateway.Authorize(ctx, order.CustomerID, totals.Grand)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrPaymentDeclined, err)
	}
	return paymentReference(order.ID, reference), nil
}

func paymentReference(orderID, gatewayReference string) string {
	return orderID + ":" + gatewayReference
}
