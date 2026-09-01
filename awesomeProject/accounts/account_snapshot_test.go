package accounts

import (
	"testing"

	"awesomeProject/accounts/models"
)

func TestAccountSnapshotReturnsCopy(t *testing.T) {
	handler := New()
	handler.accounts["alice"] = &models.Account{
		Name:   "alice",
		Amount: 50,
	}

	account, ok := handler.accountSnapshot("alice")
	if !ok {
		t.Fatal("accountSnapshot() did not find an existing account")
	}

	account.Amount = 100
	if got := handler.accounts["alice"].Amount; got != 50 {
		t.Fatalf("stored amount = %d, want 50", got)
	}
}
