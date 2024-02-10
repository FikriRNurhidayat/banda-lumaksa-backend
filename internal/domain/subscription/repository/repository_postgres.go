package subscription_repository

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/infra/logger"
	database_manager "github.com/fikrirnurhidayat/banda-lumaksa/internal/manager/database"
	"github.com/fikrirnurhidayat/banda-lumaksa/pkg/exists"
	"github.com/google/uuid"

	postgres_repository "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/repository/postgres"

	subscription_entity "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/entity"
	subscription_specification "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/specification"
	subscription_types "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/types"
)

type PostgresSubscriptionRow struct {
	ID               uuid.NullUUID
	Name             sql.NullString
	Fee              sql.NullInt32
	SubscriptionType sql.NullString
	StartedAt        sql.NullTime
	EndedAt          sql.NullTime
	DueAt            sql.NullTime
	CreatedAt        sql.NullTime
	UpdatedAt        sql.NullTime
}

var NoPostgresSubscriptionRow = PostgresSubscriptionRow{}

func NewPostgresRepository(logger logger.Logger, dbm database_manager.DatabaseManager) (SubscriptionRepository, error) {
	return postgres_repository.New[subscription_entity.Subscription, subscription_specification.SubscriptionSpecification, PostgresSubscriptionRow](postgres_repository.Option[subscription_entity.Subscription, subscription_specification.SubscriptionSpecification, PostgresSubscriptionRow]{
		Logger:    logger,
		TableName: "subscriptions",
		Schema: map[string]string{
			"id":                postgres_repository.UUID,
			"name":              postgres_repository.CharacterVarying,
			"fee":               postgres_repository.Integer,
			"subscription_type": postgres_repository.CharacterVarying,
			"started_at":        postgres_repository.TimestampWithZone,
			"ended_at":          postgres_repository.TimestampWithZone,
			"due_at":            postgres_repository.TimestampWithZone,
			"created_at":        postgres_repository.TimestampWithZone,
			"updated_at":        postgres_repository.TimestampWithZone,
		},
		Columns: []string{
			"id",
			"name",
			"fee",
			"subscription_type",
			"started_at",
			"ended_at",
			"due_at",
			"created_at",
			"updated_at",
		},
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
				ID:        row.ID.UUID,
				Name:      row.Name.String,
				Fee:       row.Fee.Int32,
				Type:      subscription_types.GetType(row.SubscriptionType.String),
				StartedAt: row.StartedAt.Time,
				EndedAt:   row.EndedAt.Time,
				DueAt:     row.DueAt.Time,
				CreatedAt: row.CreatedAt.Time,
				UpdatedAt: row.UpdatedAt.Time,
			}
		},
		Row: func(subscription subscription_entity.Subscription) PostgresSubscriptionRow {
			subscriptionType := subscription.Type.String()

			return PostgresSubscriptionRow{
				ID: uuid.NullUUID{
					UUID:  subscription.ID,
					Valid: true,
				},
				Name: sql.NullString{
					String: subscription.Name,
					Valid:  exists.String(subscription.Name),
				},
				Fee: sql.NullInt32{
					Int32: subscription.Fee,
					Valid: exists.Number(uint32(subscription.Fee)),
				},
				SubscriptionType: sql.NullString{
					String: subscriptionType,
					Valid:  subscriptionType != "",
				},
				StartedAt: sql.NullTime{
					Time:  subscription.StartedAt,
					Valid: exists.Date(subscription.StartedAt),
				},
				EndedAt: sql.NullTime{
					Time:  subscription.EndedAt,
					Valid: exists.Date(subscription.EndedAt),
				},
				DueAt: sql.NullTime{
					Time:  subscription.DueAt,
					Valid: exists.Date(subscription.DueAt),
				},
				CreatedAt: sql.NullTime{
					Time:  subscription.CreatedAt,
					Valid: exists.Date(subscription.CreatedAt),
				},
				UpdatedAt: sql.NullTime{
					Time:  subscription.UpdatedAt,
					Valid: exists.Date(subscription.UpdatedAt),
				},
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
