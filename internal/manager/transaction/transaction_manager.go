package transaction_manager

import (
	"context"
	"database/sql"
	
	manager_values "github.com/fikrirnurhidayat/banda-lumaksa/internal/manager/values"
)

type TransactionManager interface {
	Execute(ctx context.Context, fn func(context.Context) error) error
}

type TransactionManagerImpl struct {
	db *sql.DB
}

func (m *TransactionManagerImpl) Execute(ctx context.Context, fn func(context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	if err := fn(context.WithValue(ctx, manager_values.TxKey{}, tx)); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

func New(db *sql.DB) TransactionManager {
	return &TransactionManagerImpl{
		db: db,
	}
}
