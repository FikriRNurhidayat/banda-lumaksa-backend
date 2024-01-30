package subscription

import (
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/manager"
	pgrepo "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/repository/postgres"
	"github.com/google/uuid"
)

type PostgresSubscriptionRow struct {
	ID               uuid.UUID
	Name             string
	Fee              int32
	SubscriptionType string
	StartedAt        time.Time
	EndedAt          sql.NullTime
	DueAt            time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

var NoPostgresSubscriptionRow = PostgresSubscriptionRow{}

func NewPostgresSubscriptionRepository(dbm manager.DatabaseManager) SubscriptionRepository {
	return pgrepo.New[Subscription, SubscriptionSpecification, PostgresSubscriptionRow](pgrepo.Option[Subscription, SubscriptionSpecification, PostgresSubscriptionRow]{
		TableName:       "subscriptions",
		Columns:         []string{"id", "name", "fee", "subscription_type", "started_at", "ended_at", "due_at", "created_at", "updated_at"},
		PrimaryKey:      "id",
		DatabaseManager: dbm,
		Filter: func(specs ...SubscriptionSpecification) squirrel.Sqlizer {
			where := squirrel.And{}
			for _, spec := range specs {
				switch v := spec.(type) {
				case WithIDSpecification:
					where = append(where, squirrel.Eq{"id": v.ID})
				case NameLikeSpecification:
					where = append(where, squirrel.ILike{"name": v.Substring})
				case TypeIsSpecification:
					where = append(where, squirrel.Eq{"subscription_type": v.Type.String()})
				case CreatedBetweenSpecification:
					where = append(where, squirrel.LtOrEq{"created_at": v.End})
				case StartedBetweenSpecification:
					where = append(where, squirrel.LtOrEq{"started_at": v.End})
				case EndedBetweenSpecification:
					where = append(where, squirrel.LtOrEq{"ended_at": v.End})
				case NotEndedSpecification:
					where = append(where, squirrel.Or{squirrel.GtOrEq{"ended_at": v.Now}, squirrel.Eq{"ended_at": nil}})
				case DueBetweenSpecification:
					where = append(where, squirrel.LtOrEq{"due_at": v.End})
				case DueBeforeSpecification:
					where = append(where, squirrel.LtOrEq{"due_at": v.Now})
				}
			}
			return where
		},
		Scan: func(rows *sql.Rows) (PostgresSubscriptionRow, error) {
			row := PostgresSubscriptionRow{}
			if err := rows.Scan(&row.ID, &row.Name, &row.Fee, &row.SubscriptionType, &row.StartedAt, &row.EndedAt, &row.DueAt, &row.CreatedAt, &row.UpdatedAt); err != nil {
				return NoPostgresSubscriptionRow, err
			}

			return row, nil
		},
		Entity: func(row PostgresSubscriptionRow) Subscription {
			return Subscription{
				ID:        row.ID,
				Name:      row.Name,
				Fee:       row.Fee,
				Type:      GetType(row.SubscriptionType),
				StartedAt: row.StartedAt,
				EndedAt:   row.EndedAt.Time,
				DueAt:     row.DueAt,
				CreatedAt: row.CreatedAt,
				UpdatedAt: row.UpdatedAt,
			}
		},
		Row: func(subscription Subscription) PostgresSubscriptionRow {
			return PostgresSubscriptionRow{
				ID:               subscription.ID,
				Name:             subscription.Name,
				Fee:              subscription.Fee,
				SubscriptionType: subscription.Type.String(),
				StartedAt:        subscription.StartedAt,
				EndedAt: sql.NullTime{
					Time:  subscription.EndedAt,
					Valid: !subscription.EndedAt.IsZero(),
				},
				DueAt:     subscription.DueAt,
				CreatedAt: subscription.CreatedAt,
				UpdatedAt: subscription.UpdatedAt,
			}
		},
		Values: func(row PostgresSubscriptionRow) []any {
			return []any{
				row.ID,
				row.Name,
				row.Fee,
				row.SubscriptionType,
				row.StartedAt,
				row.EndedAt,
				row.DueAt,
				row.CreatedAt,
				row.UpdatedAt,
			}
		},
	})
}
