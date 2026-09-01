package reviewpathdemo

import "strings"

func normalizeOrder(order Order) Order {
	normalized := order
	normalized.ID = normalizeIdentifier(order.ID)
	normalized.CustomerID = normalizeIdentifier(order.CustomerID)
	normalized.DestinationCountry = normalizeCountry(order.DestinationCountry)
	normalized.Coupon = normalizeIdentifier(order.Coupon)
	normalized.Items = make([]Item, len(order.Items))

	for index, item := range order.Items {
		normalized.Items[index] = normalizeItem(item)
	}

	return normalized
}

func normalizeItem(item Item) Item {
	item.SKU = normalizeIdentifier(item.SKU)
	return item
}

func normalizeIdentifier(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func normalizeCountry(value string) string {
	return normalizeIdentifier(value)
}
