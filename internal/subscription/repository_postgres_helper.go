package subscription

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/repository"
)

func (*PostgresSubscriptionRepository) scan(rows *sql.Rows) (*PostgresSubscriptionRow, error) {
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

func (r *PostgresSubscriptionRepository) filter(specs ...SubscriptionSpecification) squirrel.Sqlizer {
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
}

func (r *PostgresSubscriptionRepository) query(ctx context.Context, args repository.ListArgs[SubscriptionSpecification]) (*sql.Rows, error) {
	builder := squirrel.
		Select(Columns...).
		From(TableName).
		Where(r.filter(args.Filters...))
	builder = r.dbm.Paginate(builder, args.Limit, args.Offset)
	queryStr, queryArgs, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	return r.dbm.Querier(ctx).QueryContext(ctx, queryStr, queryArgs...)
}
