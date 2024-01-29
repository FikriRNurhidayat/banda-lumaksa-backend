package manager

import (
	"context"
	"database/sql"
)

type txKey struct{}

type Querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type DatabaseManager interface {
	Querier(ctx context.Context) Querier
}

type DatabaseManagerImpl struct {
	db *sql.DB
}

func (man *DatabaseManagerImpl) Querier(ctx context.Context) Querier {
	hasExternalTransaction := ctx.Value(txKey{}) != nil
	if !hasExternalTransaction {
		return man.db
	}

	v := ctx.Value(txKey{})
	tx, ok := v.(*sql.Tx)
	if ok {
		return tx
	}

	return man.db
}

func NewDatabaseManager(db *sql.DB) DatabaseManager {
	return &DatabaseManagerImpl{
		db: db,
	}
}
