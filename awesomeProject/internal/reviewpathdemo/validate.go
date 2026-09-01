package reviewpathdemo

import "fmt"

func validateOrder(order Order) error {
	if err := validateOrderID(order.ID); err != nil {
		return err
	}
	if err := validateCustomerID(order.CustomerID); err != nil {
		return err
	}
	if err := validateDestination(order.DestinationCountry); err != nil {
		return err
	}
	return validateItems(order.Items)
}

func validateOrderID(orderID string) error {
	if orderID == "" {
		return fmt.Errorf("%w: empty order ID", ErrInvalidOrder)
	}
	return nil
}

func validateCustomerID(customerID string) error {
	if customerID == "" {
		return fmt.Errorf("%w: empty customer ID", ErrInvalidOrder)
	}
	return nil
}

func validateDestination(country string) error {
	if len(country) != 2 {
		return fmt.Errorf("%w: destination must be a two-letter country", ErrInvalidOrder)
	}
	return nil
}

func validateItems(items []Item) error {
	if len(items) == 0 {
		return fmt.Errorf("%w: order has no items", ErrInvalidOrder)
	}
	for index, item := range items {
		if err := validateItem(item); err != nil {
			return fmt.Errorf("item %d: %w", index, err)
		}
	}
	return nil
}

func validateItem(item Item) error {
	if item.SKU == "" {
		return fmt.Errorf("%w: empty SKU", ErrInvalidOrder)
	}
	if item.Quantity <= 0 {
		return fmt.Errorf("%w: quantity must be positive", ErrInvalidOrder)
	}
	if item.PriceCents < 0 {
		return fmt.Errorf("%w: price must not be negative", ErrInvalidOrder)
	}
	if item.WeightGrams < 0 {
		return fmt.Errorf("%w: weight must not be negative", ErrInvalidOrder)
	}
	return nil
}
