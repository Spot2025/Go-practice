package reviewpathdemo

func buildAuditTrail(order Order, totals Totals, payment string) []string {
	return []string{
		auditEntry("order", order.ID),
		auditEntry("customer", order.CustomerID),
		auditEntry("subtotal", formatCents(totals.Subtotal)),
		auditEntry("grand", formatCents(totals.Grand)),
		auditEntry("payment", payment),
	}
}

func auditEntry(key, value string) string {
	return key + "=" + value
}
