package accounts

import (
	"context"
	"errors"
	"testing"
)

func TestTransferSuccess(t *testing.T) {
	handler, repository, publisher := newTestTransferHandler(map[string]int64{
		"alice": 1000,
		"bob":   100,
	})

	result, err := handler.HandleTransfer(context.Background(), TransferRequest{
		SenderID:   " alice ",
		ReceiverID: " bob ",
		Amount:     100,
	})
	if err != nil {
		t.Fatalf("HandleTransfer() error = %v", err)
	}

	if result.Fee != 2 {
		t.Errorf("result fee = %d, want 2", result.Fee)
	}
	if result.SenderBalance != 898 {
		t.Errorf("sender balance = %d, want 898", result.SenderBalance)
	}
	if result.ReceiverBalance != 200 {
		t.Errorf("receiver balance = %d, want 200", result.ReceiverBalance)
	}

	senderBalance, _ := repository.Balance("alice")
	receiverBalance, _ := repository.Balance("bob")
	if senderBalance != 898 || receiverBalance != 200 {
		t.Errorf("stored balances = (%d, %d), want (898, 200)", senderBalance, receiverBalance)
	}

	events := publisher.Events()
	if len(events) != 1 {
		t.Fatalf("published events = %d, want 1", len(events))
	}
	if events[0].SenderID != "alice" || events[0].ReceiverID != "bob" {
		t.Errorf("event accounts = (%q, %q), want (alice, bob)", events[0].SenderID, events[0].ReceiverID)
	}
	if events[0].Amount != 100 || events[0].Fee != 2 {
		t.Errorf("event amount and fee = (%d, %d), want (100, 2)", events[0].Amount, events[0].Fee)
	}
	if events[0].OccurredAt.IsZero() {
		t.Error("event occurrence time is zero")
	}
}

func TestTransferInsufficientFunds(t *testing.T) {
	handler, repository, publisher := newTestTransferHandler(map[string]int64{
		"alice": 101,
		"bob":   100,
	})

	_, err := handler.HandleTransfer(context.Background(), TransferRequest{
		SenderID:   "alice",
		ReceiverID: "bob",
		Amount:     100,
	})
	if !errors.Is(err, ErrInsufficientFunds) {
		t.Fatalf("HandleTransfer() error = %v, want %v", err, ErrInsufficientFunds)
	}

	senderBalance, _ := repository.Balance("alice")
	receiverBalance, _ := repository.Balance("bob")
	if senderBalance != 101 || receiverBalance != 100 {
		t.Errorf("stored balances = (%d, %d), want (101, 100)", senderBalance, receiverBalance)
	}
	if events := publisher.Events(); len(events) != 0 {
		t.Errorf("published events = %d, want 0", len(events))
	}
}

func TestTransferRejectsNegativeAmount(t *testing.T) {
	handler, repository, publisher := newTestTransferHandler(map[string]int64{
		"alice": 100,
		"bob":   100,
	})

	_, err := handler.HandleTransfer(context.Background(), TransferRequest{
		SenderID:   "alice",
		ReceiverID: "bob",
		Amount:     -10,
	})
	if !errors.Is(err, ErrInvalidTransferAmount) {
		t.Fatalf("HandleTransfer() error = %v, want %v", err, ErrInvalidTransferAmount)
	}

	assertTransferNotApplied(t, repository, publisher, 100, 100)
}

func TestTransferRejectsSameAccount(t *testing.T) {
	handler, repository, publisher := newTestTransferHandler(map[string]int64{
		"alice": 100,
		"bob":   100,
	})

	_, err := handler.HandleTransfer(context.Background(), TransferRequest{
		SenderID:   " alice ",
		ReceiverID: "alice",
		Amount:     10,
	})
	if !errors.Is(err, ErrSameTransferAccount) {
		t.Fatalf("HandleTransfer() error = %v, want %v", err, ErrSameTransferAccount)
	}

	assertTransferNotApplied(t, repository, publisher, 100, 100)
}

func TestCalculateTransferFee(t *testing.T) {
	tests := []struct {
		name   string
		amount int64
		want   int64
	}{
		{name: "zero", amount: 0, want: 0},
		{name: "minimum fee", amount: 1, want: 1},
		{name: "exact percent", amount: 100, want: 2},
		{name: "round up fractional fee", amount: 101, want: 3},
		{name: "larger exact fee", amount: 250, want: 5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := calculateTransferFee(test.amount); got != test.want {
				t.Errorf("calculateTransferFee(%d) = %d, want %d", test.amount, got, test.want)
			}
		})
	}
}

func newTestTransferHandler(
	balances map[string]int64,
) (*TransferHandler, *MemoryTransferRepository, *MemoryTransferPublisher) {
	repository := NewMemoryTransferRepository(balances)
	publisher := &MemoryTransferPublisher{}

	return NewTransferHandler(repository, publisher), repository, publisher
}

func assertTransferNotApplied(
	t *testing.T,
	repository *MemoryTransferRepository,
	publisher *MemoryTransferPublisher,
	wantSenderBalance int64,
	wantReceiverBalance int64,
) {
	t.Helper()

	senderBalance, _ := repository.Balance("alice")
	receiverBalance, _ := repository.Balance("bob")
	if senderBalance != wantSenderBalance || receiverBalance != wantReceiverBalance {
		t.Errorf(
			"stored balances = (%d, %d), want (%d, %d)",
			senderBalance,
			receiverBalance,
			wantSenderBalance,
			wantReceiverBalance,
		)
	}
	if events := publisher.Events(); len(events) != 0 {
		t.Errorf("published events = %d, want 0", len(events))
	}
}
