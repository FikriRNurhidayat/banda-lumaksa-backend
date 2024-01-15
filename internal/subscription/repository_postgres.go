package subscription

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
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

type PostgresRepository struct {
	db *sql.DB
}

const TableName = "subscriptions"

var Columns []string = []string{
	"id",
	"name",
	"fee",
	"subscription_type",
	"started_at",
	"ended_at",
	"due_at",
	"created_at",
	"updated_at",
}

func (r *PostgresRepository) List(ctx context.Context, specs ...Specification) ([]Subscription, error) {
	var subs []Subscription
	var err error

	builder := sq.Select(Columns...).From("subscriptions")
	builder = r.Filter(builder, specs...)
	builder = r.Paginate(builder, specs...)
	queryStr, queryArgs, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return NoSubscriptions, err
	}

	rows, err := r.db.QueryContext(ctx, queryStr, queryArgs...)
	if err != nil {
		return NoSubscriptions, err
	}

	for rows.Next() {
		row, err := r.Scan(rows)
		if err != nil {
			return NoSubscriptions, err
		}

		subs = append(subs, row.Subscription())
	}

	return subs, nil
}

func (r *PostgresRepository) Size(ctx context.Context, specs ...Specification) (uint32, error) {
	var count uint32
	var err error
	builder := r.Filter(sq.Select("COUNT(id)").From("subscriptions"), specs...)
	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
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

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := sq.
		Delete("subscriptions").
		Where(sq.Eq{"id": id.String()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) Get(ctx context.Context, id uuid.UUID) (Subscription, error) {
	var subscription Subscription
	var err error

	query, args, err := sq.
		Select(Columns...).
		From("subscriptions").
		Where(sq.Eq{"id": id.String()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return NoSubscription, err
	}

	for rows.Next() {
		row, err := r.Scan(rows)
		if err != nil {
			return NoSubscription, err
		}

		subscription = row.Subscription()
	}

	return subscription, nil
}

func (r *PostgresRepository) Save(ctx context.Context, subscription Subscription) error {
	row := r.PostgresSubscriptionRow(subscription)

	query, args, err := sq.
		Insert("subscriptions").
		Columns(Columns...).
		Values(row.ID.String(), row.Name, row.Fee, row.SubscriptionType, row.StartedAt, row.EndedAt, row.DueAt, row.CreatedAt, row.UpdatedAt).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &PostgresRepository{
		db: db,
	}
}

func (*PostgresRepository) PostgresSubscriptionRow(subscription Subscription) *PostgresSubscriptionRow {
	return &PostgresSubscriptionRow{
		ID:               subscription.ID,
		Name:             subscription.Name,
		Fee:              subscription.Fee,
		SubscriptionType: subscription.Type.String(),
		StartedAt:        subscription.StartedAt,
		EndedAt: sql.NullTime{
			Time:  subscription.EndedAt,
			Valid: subscription.EndedAt.IsZero(),
		},
		DueAt:     subscription.DueAt,
		CreatedAt: subscription.CreatedAt,
		UpdatedAt: subscription.UpdatedAt,
	}
}

func (*PostgresRepository) Scan(rows *sql.Rows) (*PostgresSubscriptionRow, error) {
	row := &PostgresSubscriptionRow{}

	if err := rows.Scan(
		&row.ID,
		&row.Name,
		&row.Fee,
		&row.SubscriptionType,
		&row.StartedAt,
		&row.EndedAt,
		&row.DueAt,
		&row.CreatedAt,
		&row.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return row, nil
}

func (r *PostgresRepository) Filter(builder sq.SelectBuilder, specs ...Specification) sq.SelectBuilder {
	for _, spec := range specs {
		switch v := spec.(type) {
		case NameLikeSpecification:
			builder = builder.Where(sq.ILike{"name": v.Substring})
		case TypeIsSpecification:
			builder = builder.Where(sq.Eq{"subscription_type": v.Type.String()})
		case CreatedBetweenSpecification:
			builder = builder.Where(sq.LtOrEq{"created_at": v.End}).Where(sq.GtOrEq{"created_at": v.Start})
		case StartedBetweenSpecification:
			builder = builder.Where(sq.LtOrEq{"started_at": v.End}).Where(sq.GtOrEq{"started_at": v.Start})
		case EndedBetweenSpecification:
			builder = builder.Where(sq.LtOrEq{"ended_at": v.End}).Where(sq.GtOrEq{"ended_at": v.Start})
		case DueBetweenSpecification:
			builder = builder.Where(sq.LtOrEq{"due_at": v.End}).Where(sq.GtOrEq{"due_at": v.Start})
		}
	}

	return builder
}

func (r *PostgresRepository) Paginate(builder sq.SelectBuilder, specs ...Specification) sq.SelectBuilder {
	for _, spec := range specs {
		switch v := spec.(type) {
		case LimitSpecification:
			builder = builder.Limit(uint64(v.Limit))
		case OffsetSpecification:
			builder = builder.Offset(uint64(v.Offset))
		}
	}

	return builder
}

func (row *PostgresSubscriptionRow) Subscription() Subscription {
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
}
