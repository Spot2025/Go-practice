package reviewpathdemo

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidOrder          = errors.New("invalid order")
	ErrInsufficientInventory = errors.New("insufficient inventory")
	ErrPaymentDeclined       = errors.New("payment declined")
)

type CheckoutError struct {
	Stage string
	Err   error
}

func (e CheckoutError) Error() string {
	return fmt.Sprintf("checkout %s: %v", e.Stage, e.Err)
}

func (e CheckoutError) Unwrap() error {
	return e.Err
}
