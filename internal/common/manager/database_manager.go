package manager

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/specification"
)

type txKey struct{}

type Querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type DatabaseManager interface {
	Querier(ctx context.Context) Querier
	Paginate(builder squirrel.SelectBuilder, specs ...specification.Specification) squirrel.SelectBuilder
}

type DatabaseManagerImpl struct {
	db *sql.DB
}

// Paginate implements DatabaseManager.
func (m *DatabaseManagerImpl) Paginate(builder squirrel.SelectBuilder, specs ...specification.Specification) squirrel.SelectBuilder {
	for _, spec := range specs {
		switch v := spec.(type) {
		case specification.LimitSpecification:
			builder = builder.Limit(uint64(v.Limit))
		case specification.OffsetSpecification:
			builder = builder.Offset(uint64(v.Offset))
		}
	}

	return builder
}

func (m *DatabaseManagerImpl) Querier(ctx context.Context) Querier {
	hasExternalTransaction := ctx.Value(txKey{}) != nil
	if !hasExternalTransaction {
		return m.db
	}

	v := ctx.Value(txKey{})
	tx, ok := v.(*sql.Tx)
	if ok {
		return tx
	}

	return m.db
}

func NewDatabaseManager(db *sql.DB) DatabaseManager {
	return &DatabaseManagerImpl{
		db: db,
	}
}
