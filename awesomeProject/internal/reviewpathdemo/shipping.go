package reviewpathdemo

func calculateShipping(order Order) int64 {
	weight := totalWeight(order.Items)
	return shippingBase(order.DestinationCountry) + oversizedSurcharge(weight)
}

func totalWeight(items []Item) int {
	weight := 0
	for _, item := range items {
		weight += itemWeight(item)
	}
	return weight
}

func itemWeight(item Item) int {
	return item.WeightGrams * item.Quantity
}

func shippingBase(country string) int64 {
	switch country {
	case "US":
		return 500
	case "DE", "GB":
		return 700
	default:
		return 1200
	}
}

func oversizedSurcharge(weightGrams int) int64 {
	if weightGrams > 5000 {
		return 900
	}
	return 0
}
