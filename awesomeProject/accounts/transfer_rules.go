package accounts

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

const transferFeePercent int64 = 2

var (
	ErrSenderRequired          = errors.New("sender account is required")
	ErrReceiverRequired        = errors.New("receiver account is required")
	ErrInvalidTransferAmount   = errors.New("transfer amount must be positive")
	ErrSameTransferAccount     = errors.New("sender and receiver accounts must differ")
	ErrTransferAccountNotFound = errors.New("transfer account not found")
	ErrInsufficientFunds       = errors.New("insufficient funds")
	ErrReceiverBalanceOverflow = errors.New("receiver balance overflow")
	ErrConcurrentTransfer      = errors.New("account balances changed during transfer")
)

func normalizeTransferRequest(request TransferRequest) TransferRequest {
	request.SenderID = strings.TrimSpace(request.SenderID)
	request.ReceiverID = strings.TrimSpace(request.ReceiverID)

	return request
}

func validateTransfer(request TransferRequest) error {
	switch {
	case request.SenderID == "":
		return ErrSenderRequired
	case request.ReceiverID == "":
		return ErrReceiverRequired
	case request.Amount <= 0:
		return ErrInvalidTransferAmount
	case request.SenderID == request.ReceiverID:
		return ErrSameTransferAccount
	default:
		return nil
	}
}

func calculateTransferFee(amount int64) int64 {
	feeForWholeHundreds := amount / 100 * transferFeePercent
	remainingFeeNumerator := amount % 100 * transferFeePercent

	return feeForWholeHundreds + (remainingFeeNumerator+99)/100
}

func ensureSufficientFunds(sender TransferAccount, amount, fee int64) error {
	if sender.Balance < amount || sender.Balance-amount < fee {
		return fmt.Errorf(
			"%w: account %q has %d available",
			ErrInsufficientFunds,
			sender.ID,
			sender.Balance,
		)
	}

	return nil
}

func applyTransfer(accounts TransferAccounts, amount, fee int64) (TransferAccounts, error) {
	if accounts.Receiver.Balance > math.MaxInt64-amount {
		return TransferAccounts{}, ErrReceiverBalanceOverflow
	}

	accounts.Sender.Balance -= amount + fee
	accounts.Receiver.Balance += amount

	return accounts, nil
}
