package transaction

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/manager"
	"github.com/google/uuid"
)

type PostgresTransactionRow struct {
	ID          uuid.UUID
	Description string
	Amount      int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PostgresTransactionRepository struct {
	dbm manager.DatabaseManager
}

const TableName = "transactions"

var Columns []string = []string{
	"id",
	"description",
	"amount",
	"created_at",
	"updated_at",
}

// Delete implements TransactionRepository.
func (r *PostgresTransactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := sq.
		Delete(TableName).
		Where(sq.Eq{"id": id.String()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := r.dbm.Querier(ctx).ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *PostgresTransactionRepository) Get(ctx context.Context, id uuid.UUID) (Transaction, error) {
	var transaction Transaction
	var err error

	query, args, err := sq.
		Select(Columns...).
		From(TableName).
		Where(sq.Eq{"id": id.String()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	rows, err := r.dbm.Querier(ctx).QueryContext(ctx, query, args...)
	if err != nil {
		return NoTransaction, err
	}

	for rows.Next() {
		row, err := r.Scan(rows)
		if err != nil {
			return NoTransaction, err
		}

		transaction = row.Transaction()
	}

	return transaction, nil
}

func (r *PostgresTransactionRepository) List(ctx context.Context, specs ...TransactionSpecification) (Transactions, error) {
	var transactions Transactions
	var err error

	builder := sq.Select(Columns...).From(TableName)
	builder = r.Filter(builder, specs...)
	builder = r.Paginate(builder, specs...)
	queryStr, queryArgs, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return NoTransactions, err
	}

	rows, err := r.dbm.Querier(ctx).QueryContext(ctx, queryStr, queryArgs...)
	if err != nil {
		return NoTransactions, err
	}

	for rows.Next() {
		row, err := r.Scan(rows)
		if err != nil {
			return NoTransactions, err
		}

		transactions = append(transactions, row.Transaction())
	}

	return transactions, nil
}

func (r *PostgresTransactionRepository) Save(ctx context.Context, transaction Transaction) error {
	row := r.PostgresTransactionRow(transaction)

	query, args, err := sq.
		Insert(TableName).
		Columns(Columns...).
		Values(row.ID.String(), row.Description, row.Amount, row.CreatedAt, row.UpdatedAt).
		PlaceholderFormat(sq.Dollar).
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
	builder := r.Filter(sq.Select("COUNT(id)").From(TableName), specs...)
	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
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

func (r *PostgresTransactionRepository) Filter(builder sq.SelectBuilder, specs ...TransactionSpecification) sq.SelectBuilder {
	return builder
}

func (r *PostgresTransactionRepository) Paginate(builder sq.SelectBuilder, specs ...TransactionSpecification) sq.SelectBuilder {
	return builder
}

func (*PostgresTransactionRepository) Scan(rows *sql.Rows) (*PostgresTransactionRow, error) {
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

func NewPostgresTransactionRepository(dbm manager.DatabaseManager) TransactionRepository {
	return &PostgresTransactionRepository{
		dbm: dbm,
	}
}
