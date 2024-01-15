package subscription

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

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

func (r *PostgresRepository) List(ctx context.Context, specs ...Specification) ([]Subscription, error) {
	var subs []Subscription
	var err error

	queryStr, queryArgs := r.buildQuery(PostgresListSQL, specs...)
	rows, err := r.db.QueryContext(ctx, queryStr, queryArgs...)
	if err != nil {
		return NoSubscriptions, err
	}

	for rows.Next() {
		row, err := r.scanRow(rows)
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

	queryStr, queryArgs := r.buildQuery(PostgresSizeSQL, specs...)

	rows, err := r.db.QueryContext(ctx, queryStr, queryArgs...)
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
	if _, err := r.db.ExecContext(ctx, PostgresDeleteSQL, id.String()); err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) Get(ctx context.Context, id uuid.UUID) (Subscription, error) {
	var subscription Subscription
	var err error

	rows, err := r.db.QueryContext(ctx, PostgresGetSQL, id.String())
	if err != nil {
		return NoSubscription, err
	}

	for rows.Next() {
		row, err := r.scanRow(rows)
		if err != nil {
			return NoSubscription, err
		}

		subscription = row.Subscription()
	}

	return subscription, nil
}

func (r *PostgresRepository) Save(ctx context.Context, subscription Subscription) error {
	row := r.PostgresSubscriptionRow(subscription)
	if _, err := r.db.ExecContext(ctx, PostgresSaveSQL, row.QueryArgs()...); err != nil {
		return err
	}

	return nil
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &PostgresRepository{
		db: db,
	}
}

func (s *PostgresSubscriptionRow) NamedArgs() []sql.NamedArg {
	return []sql.NamedArg{
		sql.Named("id", s.ID.String()),
		sql.Named("name", s.Name),
		sql.Named("fee", s.Fee),
		sql.Named("subscription_type", s.SubscriptionType),
		sql.Named("started_at", s.StartedAt),
		sql.Named("ended_at", s.EndedAt),
		sql.Named("due_at", s.DueAt),
		sql.Named("created_at", s.CreatedAt),
		sql.Named("updated_at", s.UpdatedAt),
	}
}

func (s *PostgresSubscriptionRow) QueryArgs() []any {
	namedArgs := s.NamedArgs()

	qargs := make([]interface{}, len(namedArgs))
	for i, namedArg := range namedArgs {
		namedArgs[i] = namedArg
	}

	return qargs
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
			Valid: subscription.EndedAt != time.Time{},
		},
		DueAt:     subscription.DueAt,
		CreatedAt: subscription.CreatedAt,
		UpdatedAt: subscription.UpdatedAt,
	}
}

func (*PostgresRepository) scanRow(rows *sql.Rows) (*PostgresSubscriptionRow, error) {
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

func (r *PostgresRepository) buildQuery(baseQuery string, specs ...Specification) (string, []any) {
	var conditions []string
	var namedArgs []sql.NamedArg

	if len(specs) == 0 {
		return baseQuery, []any{}
	}

	for _, spec := range specs {
		queryStr, queryArgs := r.specToSQL(spec)
		if queryStr == "" {
			continue
		}
		condition := " (" + queryStr + ") "
		conditions = append(conditions, condition)
		namedArgs = append(namedArgs, queryArgs...)
	}

	if len(conditions) == 0 {
		return baseQuery, []any{}
	}

	qargs := make([]interface{}, len(namedArgs))
	for i, namedArg := range namedArgs {
		namedArgs[i] = namedArg
	}

	return baseQuery + " WHERE " + r.joinConditions(conditions, "AND"), qargs
}

func (*PostgresRepository) specToSQL(spec Specification) (string, []sql.NamedArg) {
	switch v := spec.(type) {
	case NameLikeSpecification:
		return "name ILIKE @nameLike", []sql.NamedArg{sql.Named("namedLike", v.Substring)}
	case TypeIsSpecification:
		return "subscription_type = @subscriptionType", []sql.NamedArg{sql.Named("subscriptionType", v.Type)}
	case CreatedBetweenSpecification:
		return "created_at BETWEEN @createdFrom AND @createdTo", []sql.NamedArg{sql.Named("createdFrom", v.Start), sql.Named("createdTo", v.End)}
	case StartedBetweenSpecification:
		return "started_at BETWEEN @startedFrom AND @startedTo", []sql.NamedArg{sql.Named("startedFrom", v.Start), sql.Named("startedTo", v.End)}
	case EndedBetweenSpecification:
		return "ended_at BETWEEN @endedFrom AND @endedTo", []sql.NamedArg{sql.Named("endedFrom", v.Start), sql.Named("endedTo", v.End)}
	case DueBetweenSpecification:
		return "due_at BETWEEN @dueFrom AND @dueTo", []sql.NamedArg{sql.Named("dueFrom", v.Start), sql.Named("dueTo", v.End)}
	default:
		return "", []sql.NamedArg{}
	}
}

// TODO: Move me
func (*PostgresRepository) joinConditions(conditions []string, operator string) string {
	return "(" + strings.Join(conditions, " "+operator+" ") + ")"
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

const PostgresSaveSQL = `INSERT INTO subscriptions (id, name, fee, subscription_type, started_at, ended_at, due_at, created_at, updated_at) 
VALUES (@id, @name, @fee, @subscription_type, @started_at, @ended_at, @due_at, @created_at, @updated_at)
ON CONFLICT (id) DO UPDATE 
SET 
  name = EXCLUDED.name, 
  fee = EXCLUDED.fee, 
  subscription_type = EXCLUDED.subscription_type, 
  started_at = EXCLUDED.started_at, 
  ended_at = EXCLUDED.ended_at, 
  due_at = EXCLUDED.due_at, 
  updated_at = NOW();`

const PostgresGetSQL = `SELECT 
  subscriptions.id, 
  subscriptions.name, 
  subscriptions.fee, 
  subscriptions.subscription_type,
  subscriptions.started_at, 
  subscriptions.ended_at, 
  subscriptions.due_at, 
  subscriptions.created_at, 
  subscriptions.updated_at 
FROM 
  subscriptions
WHERE 
  id = $1;`

const PostgresDeleteSQL = `DELETE
FROM 
  subscriptions 
WHERE 
  id = $1;`

const PostgresListSQL = `SELECT 
  subscriptions.id, 
  subscriptions.name, 
  subscriptions.fee, 
  subscriptions.subscription_type, 
  subscriptions.started_at, 
  subscriptions.ended_at, 
  subscriptions.due_at, 
  subscriptions.created_at, 
  subscriptions.updated_at 
FROM 
  subscriptions`

const PostgresSizeSQL = `SELECT
  COUNT(subscriptions.id)
FROM subscriptions`
