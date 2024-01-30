package transaction

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/repository"
)

func (r *PostgresTransactionRepository) filter(specs ...TransactionSpecification) squirrel.Sqlizer {
	return squirrel.And{}
}

func (*PostgresTransactionRepository) scan(rows *sql.Rows) (*PostgresTransactionRow, error) {
	row := &PostgresTransactionRow{}

	if err := rows.Scan(
		&row.ID,
		&row.Description,
		&row.Amount,
		&row.CreatedAt,
		&row.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return row, nil
}

func (r *PostgresTransactionRepository) query(ctx context.Context, args repository.ListArgs[TransactionSpecification]) (*sql.Rows, error) {
	builder := squirrel.
		Select(Columns...).
		From(TableName).
		Where(r.filter(args.Filters...))
	builder = r.dbm.Paginate(builder, args.Limit, args.Offset)
	queryStr, queryArgs, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	return r.dbm.Querier(ctx).QueryContext(ctx, queryStr, queryArgs...)
}
