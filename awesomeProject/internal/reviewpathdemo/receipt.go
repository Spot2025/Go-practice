package reviewpathdemo

import "fmt"

func buildReceipt(order Order, totals Totals, payment string) Receipt {
	return Receipt{
		OrderID:          order.ID,
		PaymentReference: payment,
		Lines:            receiptLines(order.Items),
		Totals:           totals,
	}
}

func receiptLines(items []Item) []string {
	lines := make([]string, 0, len(items))
	for _, item := range items {
		lines = append(lines, receiptLine(item))
	}
	return lines
}

func receiptLine(item Item) string {
	return fmt.Sprintf("%dx %s — %s", item.Quantity, item.SKU, formatCents(lineSubtotal(item)))
}

func formatCents(cents int64) string {
	sign := ""
	if cents < 0 {
		sign = "-"
		cents = -cents
	}
	return fmt.Sprintf("%s%d.%02d", sign, cents/100, cents%100)
}
