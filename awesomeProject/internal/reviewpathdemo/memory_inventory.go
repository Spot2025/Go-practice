package reviewpathdemo

import "context"

type MemoryInventory struct {
	stock map[string]int
}

func NewMemoryInventory(stock map[string]int) *MemoryInventory {
	copyOfStock := make(map[string]int, len(stock))
	for sku, quantity := range stock {
		copyOfStock[normalizeIdentifier(sku)] = quantity
	}
	return &MemoryInventory{stock: copyOfStock}
}

func (m *MemoryInventory) Available(_ context.Context, sku string, quantity int) bool {
	return m.stock[sku] >= quantity
}

func (m *MemoryInventory) Reserve(_ context.Context, sku string, quantity int) error {
	if m.stock[sku] < quantity {
		return ErrInsufficientInventory
	}
	m.stock[sku] -= quantity
	return nil
}
