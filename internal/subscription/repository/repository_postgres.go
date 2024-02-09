package subscription_repository

import (
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/manager"
	"github.com/google/uuid"

	postgres_repository "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/repository/postgres"

	subscription_entity "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/entity"
	subscription_specification "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/specification"
	subscription_types "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/types"
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

func NewPostgresRepository(dbm manager.DatabaseManager) SubscriptionRepository {
	return postgres_repository.New[subscription_entity.Subscription, subscription_specification.SubscriptionSpecification, PostgresSubscriptionRow](postgres_repository.Option[subscription_entity.Subscription, subscription_specification.SubscriptionSpecification, PostgresSubscriptionRow]{
		TableName:       "subscriptions",
		Columns:         []string{"id", "name", "fee", "subscription_type", "started_at", "ended_at", "due_at", "created_at", "updated_at"},
		PrimaryKey:      "id",
		DatabaseManager: dbm,
		Filter: func(specs ...subscription_specification.SubscriptionSpecification) squirrel.Sqlizer {
			where := squirrel.And{}
			for _, spec := range specs {
				switch v := spec.(type) {
				case subscription_specification.WithIDSpecification:
					where = append(where, squirrel.Eq{"id": v.ID})
				case subscription_specification.NameLikeSpecification:
					where = append(where, squirrel.ILike{"name": v.Substring})
				case subscription_specification.NameIsSpecification:
					where = append(where, squirrel.Eq{"name": v.Name})
				case subscription_specification.TypeIsSpecification:
					where = append(where, squirrel.Eq{"subscription_type": v.Type.String()})
				case subscription_specification.CreatedBetweenSpecification:
					where = append(where, squirrel.LtOrEq{"created_at": v.End})
				case subscription_specification.StartedBetweenSpecification:
					where = append(where, squirrel.LtOrEq{"started_at": v.End})
				case subscription_specification.EndedBetweenSpecification:
					where = append(where, squirrel.LtOrEq{"ended_at": v.End})
				case subscription_specification.NotEndedSpecification:
					where = append(where, squirrel.Or{squirrel.GtOrEq{"ended_at": v.Now}, squirrel.Eq{"ended_at": nil}})
				case subscription_specification.DueBetweenSpecification:
					where = append(where, squirrel.LtOrEq{"due_at": v.End})
				case subscription_specification.DueBeforeSpecification:
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
		Entity: func(row PostgresSubscriptionRow) subscription_entity.Subscription {
			return subscription_entity.Subscription{
				ID:        row.ID,
				Name:      row.Name,
				Fee:       row.Fee,
				Type:      subscription_types.GetType(row.SubscriptionType),
				StartedAt: row.StartedAt,				EndedAt:   row.EndedAt.Time,
				DueAt:     row.DueAt,
				CreatedAt: row.CreatedAt,
				UpdatedAt: row.UpdatedAt,
			}
		},
		Row: func(subscription subscription_entity.Subscription) PostgresSubscriptionRow {
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
