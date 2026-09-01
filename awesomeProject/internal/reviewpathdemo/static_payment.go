package reviewpathdemo

import "context"

type StaticPaymentGateway struct {
	Reference string
	Err       error
}

func (g StaticPaymentGateway) Authorize(
	_ context.Context,
	_ string,
	_ int64,
) (string, error) {
	if g.Err != nil {
		return "", g.Err
	}
	if g.Reference == "" {
		return "approved", nil
	}
	return g.Reference, nil
}
