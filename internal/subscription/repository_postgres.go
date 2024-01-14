package subscription

import (
	"context"
	"database/sql"
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
	var query string = PostgresListSQL

	queryStr, queryArgs := r.buildQuery(specs...)

	stmt, err := r.db.PrepareContext(ctx, query+queryStr)
	if err != nil {
		return EmptySubscriptions, err
	}

	rows, err := stmt.QueryContext(ctx, queryArgs)
	if err != nil {
		return EmptySubscriptions, err
	}

	for rows.Next() {
		row, err := r.scanRow(rows)
		if err != nil {
			return EmptySubscriptions, err
		}

		subs = append(subs, row.Subscription())
	}

	return subs, nil
}

func (r *PostgresRepository) Size(ctx context.Context, specs ...Specification) (uint32, error) {
	var count uint32
	var err error
	var query string = PostgresSizeSQL

	queryStr, queryArgs := r.buildQuery(specs...)

	stmt, err := r.db.PrepareContext(ctx, query+queryStr)
	if err != nil {
		return 0, err
	}

	rows, err := stmt.QueryContext(ctx, queryArgs)
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
	var err error
	var query string = PostgresDeleteSQL

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	if _, err := stmt.QueryContext(ctx, sql.Named("id", id)); err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) Get(ctx context.Context, id uuid.UUID) (Subscription, error) {
	var subscription Subscription
	var err error
	var query string = PostgresGetSQL

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return EmptySubscription, err
	}

	rows, err := stmt.QueryContext(ctx, sql.Named("id", id))
	if err != nil {
		return EmptySubscription, err
	}

	for rows.Next() {
		row, err := r.scanRow(rows)
		if err != nil {
			return EmptySubscription, err
		}

		subscription = row.Subscription()
	}

	return subscription, nil
}

func (r *PostgresRepository) Save(ctx context.Context, subscription Subscription) error {
	query := PostgresSaveSQL
	sargs := r.PostgresSubscriptionRow(subscription)

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	if _, err := stmt.QueryContext(ctx, sargs.NamedArg()); err != nil {
		return err
	}

	return nil
}

func NewPostgresRepository() Repository {
	return &PostgresRepository{}
}

func (s *PostgresSubscriptionRow) NamedArg() []sql.NamedArg {
	return []sql.NamedArg{
		sql.Named("id", s.ID),
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

func (r *PostgresRepository) buildQuery(specs ...Specification) (string, []sql.NamedArg) {
	var conditions []string
	var args []sql.NamedArg

	if len(specs) == 0 {
		return "", []sql.NamedArg{}
	}

	for _, spec := range specs {
		queryStr, queryArgs := r.specToSQL(spec)
		condition := " (" + queryStr + ") "
		conditions = append(conditions, condition)
		args = append(args, queryArgs...)
	}

	return " WHERE " + r.joinConditions(conditions, "AND"), args
}

func (*PostgresRepository) specToSQL(spec Specification) (string, []sql.NamedArg) {
	switch v := spec.(type) {
	case NameLikeSpecification:
		return "name ILIKE :nameLike", []sql.NamedArg{sql.Named("namedLike", v.Substring)}
	case TypeIsSpecification:
		return "subscription_type = :subscriptionType", []sql.NamedArg{sql.Named("subscriptionType", v.Type)}
	case CreatedBetweenSpecification:
		return "created_at BETWEEN :createdFrom AND :createdTo", []sql.NamedArg{sql.Named("createdFrom", v.Start), sql.Named("createdTo", v.End)}
	case StartedBetweenSpecification:
		return "started_at BETWEEN :startedFrom AND :startedTo", []sql.NamedArg{sql.Named("startedFrom", v.Start), sql.Named("startedTo", v.End)}
	case EndedBetweenSpecification:
		return "ended_at BETWEEN :endedFrom AND :endedTo", []sql.NamedArg{sql.Named("endedFrom", v.Start), sql.Named("endedTo", v.End)}
	case DueBetweenSpecification:
		return "due_at BETWEEN :dueFrom AND :dueTo", []sql.NamedArg{sql.Named("dueFrom", v.Start), sql.Named("dueTo", v.End)}
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

const PostgresSaveSQL = `
INSERT INTO subscriptions (
  id, name, fee, subscription_type, started_at, 
  ended_at, due_at, created_at, updated_at
) 
VALUES 
  (
    :id, :name, :fee, :subscription_type, 
    :started_at, :ended_at, :due_at, 
    :created_at, :updated_at
  ) ON CONFLICT (id) DO 
UPDATE 
SET 
  name = :name, 
  fee = :fee, 
  subscription_type = :subscription_type, 
  started_at = :started_at, 
  ended_at = :ended_at, 
  due_at = :due_at, 
  updated_at = NOW();
`

const PostgresGetSQL = `
SELECT 
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
  id = $1;
`

const PostgresDeleteSQL = `
DELETE
FROM 
  subscriptions 
WHERE 
  id = $1;
`

const PostgresListSQL = `
SELECT 
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
`

const PostgresSizeSQL = `
SELECT
  COUNT(subscriptions.id)
FROM subscriptions`
