package transaction

import (
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/manager"
	pgrepo "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/repository/postgres"
	"github.com/google/uuid"
)

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

func NewPostgresTransactionRepository(dbm manager.DatabaseManager) TransactionRepository {
	return pgrepo.New[Transaction, TransactionSpecification, *PostgresTransactionRow](pgrepo.Option[Transaction, TransactionSpecification, *PostgresTransactionRow]{
		TableName:       "transactions",
		Columns:         []string{"id", "description", "amount", "created_at", "updated_at"},
		PrimaryKey:      "id",
		DatabaseManager: dbm,
		Filter: func(specs ...TransactionSpecification) squirrel.Sqlizer {
			return squirrel.And{}
		},
		Scan: func(rows *sql.Rows) (*PostgresTransactionRow, error) {
			row := &PostgresTransactionRow{}
			if err := rows.Scan(&row.ID, &row.Description, &row.Amount, &row.CreatedAt, &row.UpdatedAt); err != nil {
				return nil, err
			}
			return row, nil
		},
		Entity: func(row *PostgresTransactionRow) Transaction {
			return Transaction{
				ID:          row.ID,
				Description: row.Description,
				Amount:      row.Amount,
				CreatedAt:   row.CreatedAt,
				UpdatedAt:   row.UpdatedAt,
			}
		},
		Row: func(transaction Transaction) *PostgresTransactionRow {
			return &PostgresTransactionRow{
				ID:          transaction.ID,
				Description: transaction.Description,
				Amount:      transaction.Amount,
				CreatedAt:   transaction.CreatedAt,
				UpdatedAt:   transaction.UpdatedAt,
			}
		},
		Values: func(row *PostgresTransactionRow) []any {
			return []any{
				row.ID,
				row.Description,
				row.Amount,
				row.CreatedAt,
				row.UpdatedAt,
			}
		},
	})
}
