package reviewpathdemo

func calculateSubtotal(items []Item) int64 {
	lines := make([]int64, 0, len(items))
	for _, item := range items {
		lines = append(lines, lineSubtotal(item))
	}
	return sumMoney(lines...)
}

func lineSubtotal(item Item) int64 {
	return multiplyMoney(item.PriceCents, item.Quantity)
}
