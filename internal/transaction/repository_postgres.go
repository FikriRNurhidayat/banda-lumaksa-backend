package transaction

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/manager"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/repository"
	"github.com/google/uuid"
)

const TableName = "transactions"

var Columns []string = []string{
	"id",
	"description",
	"amount",
	"created_at",
	"updated_at",
}

type PostgresTransactionRow struct {
	ID          uuid.UUID
	Description string
	Amount      int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (*PostgresTransactionRepository) PostgresTransactionRow(transaction Transaction) *PostgresTransactionRow {
	return &PostgresTransactionRow{
		ID:          transaction.ID,
		Description: transaction.Description,
		Amount:      transaction.Amount,
		CreatedAt:   transaction.CreatedAt,
		UpdatedAt:   transaction.UpdatedAt,
	}
}

func (row *PostgresTransactionRow) Transaction() Transaction {
	return Transaction{
		ID:          row.ID,
		Description: row.Description,
		Amount:      row.Amount,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func (row *PostgresTransactionRow) Values() []any {
	return []any{
		row.ID,
		row.Description,
		row.Amount,
		row.CreatedAt,
		row.UpdatedAt,
	}
}

type PostgresTransactionRepository struct {
	dbm manager.DatabaseManager
}

func (*PostgresTransactionRepository) Each(context.Context, repository.ListArgs[TransactionSpecification]) (repository.Iterator[Transaction], error) {
	panic("unimplemented")
}

func (r *PostgresTransactionRepository) Delete(ctx context.Context, specs ...TransactionSpecification) error {
	query, args, err := squirrel.
		Delete(TableName).
		Where(r.filter(specs...)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := r.dbm.Querier(ctx).ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *PostgresTransactionRepository) Get(ctx context.Context, specs ...TransactionSpecification) (Transaction, error) {
	var transaction Transaction
	var err error

	query, args, err := squirrel.
		Select(Columns...).
		From(TableName).
		Where(r.filter(specs...)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	rows, err := r.dbm.Querier(ctx).QueryContext(ctx, query, args...)
	if err != nil {
		return NoTransaction, err
	}

	for rows.Next() {
		row, err := r.scan(rows)
		if err != nil {
			return NoTransaction, err
		}
		transaction = row.Transaction()
	}

	return transaction, nil
}

func (r *PostgresTransactionRepository) List(ctx context.Context, args repository.ListArgs[TransactionSpecification]) ([]Transaction, error) {
	builder := squirrel.Select(Columns...).From(TableName).Where(r.filter(args.Filters...))
	builder = r.dbm.Paginate(builder, args.Limit, args.Offset)
	queryStr, queryArgs, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return NoTransactions, err
	}

	rows, err := r.dbm.Querier(ctx).QueryContext(ctx, queryStr, queryArgs...)
	if err != nil {
		return NoTransactions, err
	}

	transactions := []Transaction{}
	for rows.Next() {
		row, err := r.scan(rows)
		if err != nil {
			return NoTransactions, err
		}

		transactions = append(transactions, row.Transaction())
	}

	return transactions, nil
}

func (r *PostgresTransactionRepository) Save(ctx context.Context, transaction Transaction) error {
	row := r.PostgresTransactionRow(transaction)

	query, args, err := squirrel.
		Insert(TableName).
		Columns(Columns...).
		Values(row.Values()...).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := r.dbm.Querier(ctx).ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *PostgresTransactionRepository) Size(ctx context.Context, specs ...TransactionSpecification) (uint32, error) {
	var count uint32
	var err error
	builder := squirrel.Select("COUNT(id)").From(TableName).Where(r.filter(specs...))
	query, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return 0, err
	}

	rows, err := r.dbm.Querier(ctx).QueryContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
	}

	return count, nil
}

func NewPostgresTransactionRepository(dbm manager.DatabaseManager) TransactionRepository {
	return &PostgresTransactionRepository{
		dbm: dbm,
	}
}
