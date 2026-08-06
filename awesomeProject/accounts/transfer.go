package accounts

import (
	"context"
	"fmt"
	"time"
)

type TransferRequest struct {
	SenderID   string
	ReceiverID string
	Amount     int64
}

type TransferResult struct {
	Fee             int64
	SenderBalance   int64
	ReceiverBalance int64
}

type TransferEvent struct {
	SenderID   string
	ReceiverID string
	Amount     int64
	Fee        int64
	OccurredAt time.Time
}

type TransferAccount struct {
	ID      string
	Balance int64
}

type TransferAccounts struct {
	Sender   TransferAccount
	Receiver TransferAccount
}

type TransferRepository interface {
	LoadTransferAccounts(ctx context.Context, senderID, receiverID string) (TransferAccounts, error)
	SaveTransferAccounts(ctx context.Context, before, after TransferAccounts) error
}

type TransferPublisher interface {
	PublishTransferEvent(ctx context.Context, event TransferEvent) error
}

type TransferService struct {
	repository TransferRepository
	publisher  TransferPublisher
}

type TransferHandler struct {
	service *TransferService
}

func NewTransferHandler(repository TransferRepository, publisher TransferPublisher) *TransferHandler {
	return &TransferHandler{
		service: &TransferService{
			repository: repository,
			publisher:  publisher,
		},
	}
}

func (h *TransferHandler) HandleTransfer(ctx context.Context, request TransferRequest) (TransferResult, error) {
	result, err := h.service.Transfer(ctx, request)
	if err != nil {
		return TransferResult{}, fmt.Errorf("handle account transfer: %w", err)
	}

	return result, nil
}

func (s *TransferService) Transfer(ctx context.Context, request TransferRequest) (TransferResult, error) {
	request = normalizeTransferRequest(request)
	if err := validateTransfer(request); err != nil {
		return TransferResult{}, err
	}

	fee := calculateTransferFee(request.Amount)
	accounts, err := loadTransferAccounts(ctx, s.repository, request)
	if err != nil {
		return TransferResult{}, err
	}

	if err := ensureSufficientFunds(accounts.Sender, request.Amount, fee); err != nil {
		return TransferResult{}, err
	}

	updatedAccounts, err := applyTransfer(accounts, request.Amount, fee)
	if err != nil {
		return TransferResult{}, err
	}

	if err := saveTransfer(ctx, s.repository, accounts, updatedAccounts); err != nil {
		return TransferResult{}, err
	}

	event := buildTransferEvent(request, fee)
	if err := publishTransferEvent(ctx, s.publisher, event); err != nil {
		return TransferResult{}, err
	}

	return TransferResult{
		Fee:             fee,
		SenderBalance:   updatedAccounts.Sender.Balance,
		ReceiverBalance: updatedAccounts.Receiver.Balance,
	}, nil
}

func loadTransferAccounts(
	ctx context.Context,
	repository TransferRepository,
	request TransferRequest,
) (TransferAccounts, error) {
	accounts, err := repository.LoadTransferAccounts(ctx, request.SenderID, request.ReceiverID)
	if err != nil {
		return TransferAccounts{}, fmt.Errorf("load transfer accounts: %w", err)
	}

	return accounts, nil
}

func saveTransfer(
	ctx context.Context,
	repository TransferRepository,
	before TransferAccounts,
	after TransferAccounts,
) error {
	if err := repository.SaveTransferAccounts(ctx, before, after); err != nil {
		return fmt.Errorf("save transfer balances: %w", err)
	}

	return nil
}

func buildTransferEvent(request TransferRequest, fee int64) TransferEvent {
	return TransferEvent{
		SenderID:   request.SenderID,
		ReceiverID: request.ReceiverID,
		Amount:     request.Amount,
		Fee:        fee,
		OccurredAt: time.Now().UTC(),
	}
}

func publishTransferEvent(ctx context.Context, publisher TransferPublisher, event TransferEvent) error {
	if err := publisher.PublishTransferEvent(ctx, event); err != nil {
		return fmt.Errorf("publish transfer event: %w", err)
	}

	return nil
}
