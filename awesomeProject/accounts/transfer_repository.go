package accounts

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
)

type PostgresTransferRepository struct {
	db *sql.DB
}

func NewPostgresTransferRepository(db *sql.DB) *PostgresTransferRepository {
	return &PostgresTransferRepository{db: db}
}

func (r *PostgresTransferRepository) LoadTransferAccounts(
	ctx context.Context,
	senderID string,
	receiverID string,
) (TransferAccounts, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT name, amount FROM accounts WHERE name = $1 OR name = $2",
		senderID,
		receiverID,
	)
	if err != nil {
		return TransferAccounts{}, fmt.Errorf("query transfer accounts: %w", err)
	}
	defer rows.Close()

	accountsByID := make(map[string]TransferAccount, 2)
	for rows.Next() {
		var account TransferAccount
		if err := rows.Scan(&account.ID, &account.Balance); err != nil {
			return TransferAccounts{}, fmt.Errorf("scan transfer account: %w", err)
		}

		accountsByID[account.ID] = account
	}
	if err := rows.Err(); err != nil {
		return TransferAccounts{}, fmt.Errorf("iterate transfer accounts: %w", err)
	}

	sender, senderFound := accountsByID[senderID]
	if !senderFound {
		return TransferAccounts{}, fmt.Errorf("%w: %s", ErrTransferAccountNotFound, senderID)
	}

	receiver, receiverFound := accountsByID[receiverID]
	if !receiverFound {
		return TransferAccounts{}, fmt.Errorf("%w: %s", ErrTransferAccountNotFound, receiverID)
	}

	return TransferAccounts{Sender: sender, Receiver: receiver}, nil
}

func (r *PostgresTransferRepository) SaveTransferAccounts(
	ctx context.Context,
	before TransferAccounts,
	after TransferAccounts,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transfer transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	result, err := tx.ExecContext(ctx, `
		UPDATE accounts AS account
		SET amount = balances.new_balance
		FROM (VALUES
			($1::text, $2::integer, $3::integer),
			($4::text, $5::integer, $6::integer)
		) AS balances(name, old_balance, new_balance)
		WHERE account.name = balances.name
			AND account.amount = balances.old_balance`,
		before.Sender.ID,
		before.Sender.Balance,
		after.Sender.Balance,
		before.Receiver.ID,
		before.Receiver.Balance,
		after.Receiver.Balance,
	)
	if err != nil {
		return fmt.Errorf("update transfer balances: %w", err)
	}

	updatedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("count updated transfer accounts: %w", err)
	}
	if updatedRows != 2 {
		return ErrConcurrentTransfer
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transfer transaction: %w", err)
	}

	return nil
}

type MemoryTransferRepository struct {
	guard    sync.RWMutex
	balances map[string]int64
}

func NewMemoryTransferRepository(balances map[string]int64) *MemoryTransferRepository {
	storedBalances := make(map[string]int64, len(balances))
	for accountID, balance := range balances {
		storedBalances[accountID] = balance
	}

	return &MemoryTransferRepository{balances: storedBalances}
}

func (r *MemoryTransferRepository) LoadTransferAccounts(
	_ context.Context,
	senderID string,
	receiverID string,
) (TransferAccounts, error) {
	r.guard.RLock()
	defer r.guard.RUnlock()

	senderBalance, senderFound := r.balances[senderID]
	if !senderFound {
		return TransferAccounts{}, fmt.Errorf("%w: %s", ErrTransferAccountNotFound, senderID)
	}

	receiverBalance, receiverFound := r.balances[receiverID]
	if !receiverFound {
		return TransferAccounts{}, fmt.Errorf("%w: %s", ErrTransferAccountNotFound, receiverID)
	}

	return TransferAccounts{
		Sender:   TransferAccount{ID: senderID, Balance: senderBalance},
		Receiver: TransferAccount{ID: receiverID, Balance: receiverBalance},
	}, nil
}

func (r *MemoryTransferRepository) SaveTransferAccounts(
	_ context.Context,
	before TransferAccounts,
	after TransferAccounts,
) error {
	r.guard.Lock()
	defer r.guard.Unlock()

	if r.balances[before.Sender.ID] != before.Sender.Balance ||
		r.balances[before.Receiver.ID] != before.Receiver.Balance {
		return ErrConcurrentTransfer
	}

	r.balances[after.Sender.ID] = after.Sender.Balance
	r.balances[after.Receiver.ID] = after.Receiver.Balance

	return nil
}

func (r *MemoryTransferRepository) Balance(accountID string) (int64, bool) {
	r.guard.RLock()
	defer r.guard.RUnlock()

	balance, ok := r.balances[accountID]
	return balance, ok
}

type MemoryTransferPublisher struct {
	guard  sync.RWMutex
	events []TransferEvent
}

func (p *MemoryTransferPublisher) PublishTransferEvent(_ context.Context, event TransferEvent) error {
	p.guard.Lock()
	defer p.guard.Unlock()

	p.events = append(p.events, event)
	return nil
}

func (p *MemoryTransferPublisher) Events() []TransferEvent {
	p.guard.RLock()
	defer p.guard.RUnlock()

	events := make([]TransferEvent, len(p.events))
	copy(events, p.events)

	return events
}
