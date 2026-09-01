package reviewpathdemo

type Item struct {
	SKU         string
	PriceCents  int64
	Quantity    int
	WeightGrams int
}

type Order struct {
	ID                 string
	CustomerID         string
	DestinationCountry string
	Coupon             string
	Items              []Item
}

type Totals struct {
	Subtotal int64
	Discount int64
	Tax      int64
	Shipping int64
	Grand    int64
}

type Receipt struct {
	OrderID          string
	PaymentReference string
	Lines            []string
	Totals           Totals
}

type Result struct {
	Receipt Receipt
	Audit   []string
}
