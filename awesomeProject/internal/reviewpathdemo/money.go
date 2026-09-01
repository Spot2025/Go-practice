package reviewpathdemo

func sumMoney(values ...int64) int64 {
	var total int64
	for _, value := range values {
		total += value
	}
	return total
}

func multiplyMoney(unitPrice int64, quantity int) int64 {
	return unitPrice * int64(quantity)
}

func percentage(value int64, rate int64) int64 {
	return value * rate / 100
}

func clampMoney(value, minimum, maximum int64) int64 {
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}
