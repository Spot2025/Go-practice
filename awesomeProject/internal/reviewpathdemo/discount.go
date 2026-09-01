package reviewpathdemo

func calculateDiscount(order Order, subtotal int64) int64 {
	discount := percentage(subtotal, couponRate(order.Coupon))
	return clampMoney(discount, 0, subtotal)
}

func couponRate(coupon string) int64 {
	switch coupon {
	case "REVIEW10":
		return 10
	case "REVIEW20":
		return 20
	default:
		return 0
	}
}
