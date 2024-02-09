package database_manager

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	
	common_specification "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/specification"
	manager_values "github.com/fikrirnurhidayat/banda-lumaksa/internal/manager/values"
)

type Querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type DatabaseManager interface {
	Querier(ctx context.Context) Querier
	Paginate(builder squirrel.SelectBuilder, specs ...common_specification.Specification) squirrel.SelectBuilder
}

type DatabaseManagerImpl struct {
	db *sql.DB
}

func (m *DatabaseManagerImpl) Paginate(builder squirrel.SelectBuilder, specs ...common_specification.Specification) squirrel.SelectBuilder {
	for _, spec := range specs {
		switch v := spec.(type) {
		case common_specification.LimitSpecification:
			builder = builder.Limit(uint64(v.Limit))
		case common_specification.OffsetSpecification:
			builder = builder.Offset(uint64(v.Offset))
		}
	}

	return builder
}

func (m *DatabaseManagerImpl) Querier(ctx context.Context) Querier {
	hasExternalTransaction := ctx.Value(manager_values.TxKey{}) != nil
	if !hasExternalTransaction {
		return m.db
	}

	v := ctx.Value(manager_values.TxKey{})
	tx, ok := v.(*sql.Tx)
	if ok {
		return tx
	}

	return m.db
}

func New(db *sql.DB) DatabaseManager {
	return &DatabaseManagerImpl{
		db: db,
	}
}
