package reviewpathdemo

func calculateTax(country string, taxable int64) int64 {
	return percentage(taxable, countryTaxRate(country))
}

func taxableAmount(subtotal, discount int64) int64 {
	return clampMoney(subtotal-discount, 0, subtotal)
}

func countryTaxRate(country string) int64 {
	switch country {
	case "US":
		return 8
	case "DE":
		return 19
	case "GB":
		return 20
	default:
		return 0
	}
}
