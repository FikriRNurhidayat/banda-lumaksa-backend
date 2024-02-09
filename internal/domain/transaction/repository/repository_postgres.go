package transaction_repository

import (
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/infra/logger"

	postgres_repository "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/repository/postgres"
	database_manager "github.com/fikrirnurhidayat/banda-lumaksa/internal/manager/database"

	transaction_entity "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/transaction/entity"
	transaction_specification "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/transaction/specification"

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

func NewPostgresRepository(logger logger.Logger, dbm database_manager.DatabaseManager) TransactionRepository {
	return postgres_repository.New[transaction_entity.Transaction, transaction_specification.TransactionSpecification, *PostgresTransactionRow](postgres_repository.Option[transaction_entity.Transaction, transaction_specification.TransactionSpecification, *PostgresTransactionRow]{
		Logger:    logger,
		TableName: "transactions",
		Schema: map[string]string{
			"id":          postgres_repository.UUID,
			"description": postgres_repository.CharacterVarying,
			"amount":      postgres_repository.Integer,
			"created_at":  postgres_repository.TimestampWithZone,
			"updated_at":  postgres_repository.TimestampWithZone,
		},
		Columns: []string{
			"id",
			"description",
			"amount",
			"created_at",
			"updated_at",
		},
		PrimaryKey:      "id",
		DatabaseManager: dbm,
		Filter: func(specs ...transaction_specification.TransactionSpecification) squirrel.Sqlizer {
			return squirrel.And{}
		},
		Scan: func(rows *sql.Rows) (*PostgresTransactionRow, error) {
			row := &PostgresTransactionRow{}
			if err := rows.Scan(&row.ID, &row.Description, &row.Amount, &row.CreatedAt, &row.UpdatedAt); err != nil {
				return nil, err
			}
			return row, nil
		},
		Entity: func(row *PostgresTransactionRow) transaction_entity.Transaction {
			return transaction_entity.Transaction{
				ID:          row.ID,
				Description: row.Description,
				Amount:      row.Amount,
				CreatedAt:   row.CreatedAt,
				UpdatedAt:   row.UpdatedAt,
			}
		},
		Row: func(transaction transaction_entity.Transaction) *PostgresTransactionRow {
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
