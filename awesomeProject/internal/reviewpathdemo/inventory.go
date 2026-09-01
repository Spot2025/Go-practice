package reviewpathdemo

import (
	"context"
	"fmt"
)

type Inventory interface {
	Available(ctx context.Context, sku string, quantity int) bool
	Reserve(ctx context.Context, sku string, quantity int) error
}

func reserveOrder(ctx context.Context, inventory Inventory, order Order) error {
	for _, item := range order.Items {
		if err := reserveItem(ctx, inventory, item); err != nil {
			return err
		}
	}
	return nil
}

func reserveItem(ctx context.Context, inventory Inventory, item Item) error {
	if !inventory.Available(ctx, item.SKU, item.Quantity) {
		return fmt.Errorf("%w: %s", ErrInsufficientInventory, item.SKU)
	}
	if err := inventory.Reserve(ctx, item.SKU, item.Quantity); err != nil {
		return fmt.Errorf("reserve %s: %w", item.SKU, err)
	}
	return nil
}
