package subscription

import (
	"context"
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/manager"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/repository"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/specification"
	"github.com/google/uuid"
)

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

func (*PostgresSubscriptionRepository) PostgresSubscriptionRow(subscription Subscription) *PostgresSubscriptionRow {
	return &PostgresSubscriptionRow{
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

func (row *PostgresSubscriptionRow) Values() []any {
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
}

type PostgresSubscriptionRepository struct {
	dbm manager.DatabaseManager
}

type PostgresSubscriptionIterator struct {
	current    *sql.Rows
	rows       *sql.Rows
	repository *PostgresSubscriptionRepository
}

func (i *PostgresSubscriptionIterator) Entry() (Subscription, error) {
	row, err := i.repository.scan(i.rows)
	if err != nil {
		return NoSubscription, err
	}

	return row.Subscription(), nil
}

func (i *PostgresSubscriptionIterator) Next() bool {
	next := i.rows.Next()
	if next {
		i.current = i.rows
		return true
	}

	return false
}

func (r *PostgresSubscriptionRepository) Each(ctx context.Context, args repository.ListArgs[SubscriptionSpecification]) (repository.Iterator[Subscription], error) {
	rows, err := r.query(ctx, args)
	if err != nil {
		return nil, err
	}

	return &PostgresSubscriptionIterator{
		current:    rows,
		rows:       rows,
		repository: r,
	}, nil
}

func (r *PostgresSubscriptionRepository) List(ctx context.Context, args repository.ListArgs[SubscriptionSpecification]) ([]Subscription, error) {
	rows, err := r.query(ctx, args)
	if err != nil {
		return NoSubscriptions, err
	}

	subs := []Subscription{}
	for rows.Next() {
		row, err := r.scan(rows)
		if err != nil {
			return NoSubscriptions, err
		}

		subs = append(subs, row.Subscription())
	}

	return subs, nil
}

func (r *PostgresSubscriptionRepository) Size(ctx context.Context, specs ...SubscriptionSpecification) (uint32, error) {
	var count uint32
	var err error
	builder := squirrel.
		Select("COUNT(id)").
		From(TableName).
		Where(r.filter(specs...))
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

func (r *PostgresSubscriptionRepository) Delete(ctx context.Context, specs ...SubscriptionSpecification) error {
	builder := squirrel.Select(Columns...).From(TableName)
	query, args, err := builder.
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

func (r *PostgresSubscriptionRepository) Get(ctx context.Context, specs ...SubscriptionSpecification) (Subscription, error) {
	rows, err := r.query(ctx, repository.ListArgs[SubscriptionSpecification]{
		Filters: specs,
		Limit:   specification.WithLimit(1),
	})
	if err != nil {
		return NoSubscription, err
	}

	for rows.Next() {
		row, err := r.scan(rows)
		if err != nil {
			return NoSubscription, err
		}

		return row.Subscription(), nil
	}

	return NoSubscription, nil
}

func (r *PostgresSubscriptionRepository) Save(ctx context.Context, subscription Subscription) error {
	row := r.PostgresSubscriptionRow(subscription)

	query, args, err := squirrel.
		Insert(TableName).
		Columns(Columns...).
		Values(row.Values()...).
		Suffix("ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, fee = EXCLUDED.fee, subscription_type = EXCLUDED.subscription_type, started_at = EXCLUDED.started_at, ended_at = EXCLUDED.ended_at, due_at = EXCLUDED.due_at, updated_at = EXCLUDED.updated_at").
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

func NewPostgresSubscriptionRepository(dbm manager.DatabaseManager) SubscriptionRepository {
	return &PostgresSubscriptionRepository{
		dbm: dbm,
	}
}
