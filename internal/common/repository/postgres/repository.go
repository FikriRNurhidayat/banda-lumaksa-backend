package postgres_repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/manager"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/repository"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/specification"
)

type PostgresRepository[Entity any, Specification any, Row any] struct {
	dbm          manager.DatabaseManager
	tableName    string
	columns      []string
	primaryKey   string
	filter       func(...Specification) squirrel.Sqlizer
	scan         func(*sql.Rows) (Row, error)
	row          func(Entity) Row
	values       func(Row) []any
	entity       func(Row) Entity
	noEntities   []Entity
	noEntity     Entity
	noRow        Row
	noRows       []Row
	upsertSuffix string
}

type PostgresIterator[Entity any, Row any] struct {
	rows     *sql.Rows
	scan     func(*sql.Rows) (Row, error)
	entity   func(Row) Entity
	noEntity Entity
}

type Option[Entity any, Specification any, Row any] struct {
	TableName       string
	Columns         []string
	PrimaryKey      string
	DatabaseManager manager.DatabaseManager
	Filter          func(...Specification) squirrel.Sqlizer
	Scan            func(rows *sql.Rows) (Row, error)
	Entity          func(Row) Entity
	Row             func(Entity) Row
	Values          func(Row) []any
}

func (i *PostgresIterator[Entity, Row]) Entry() (Entity, error) {
	row, err := i.scan(i.rows)
	if err != nil {
		return i.noEntity, err
	}

	return i.entity(row), nil
}

func (i *PostgresIterator[Entity, Row]) Next() bool {
	return i.rows.Next()
}

func (r *PostgresRepository[Entity, Specification, Row]) Delete(ctx context.Context, specs ...Specification) error {
	query, args, err := squirrel.
		Delete(r.tableName).
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

func (r *PostgresRepository[Entity, Specification, Row]) Each(ctx context.Context, args repository.ListArgs[Specification]) (repository.Iterator[Entity], error) {
	rows, err := r.query(ctx, args)
	if err != nil {
		return nil, err
	}

	return &PostgresIterator[Entity, Row]{
		rows:     rows,
		scan:     r.scan,
		entity:   r.entity,
		noEntity: r.noEntity,
	}, nil
}

func (r *PostgresRepository[Entity, Specification, Row]) Get(ctx context.Context, specs ...Specification) (Entity, error) {
	rows, err := r.query(ctx, repository.ListArgs[Specification]{
		Filters: specs,
		Limit:   specification.WithLimit(1),
	})
	if err != nil {
		return r.noEntity, err
	}

	for rows.Next() {
		row, err := r.scan(rows)
		if err != nil {
			return r.noEntity, err
		}

		return r.entity(row), nil
	}

	return r.noEntity, nil
}

func (r *PostgresRepository[Entity, Specification, Row]) List(ctx context.Context, args repository.ListArgs[Specification]) ([]Entity, error) {
	rows, err := r.query(ctx, args)
	if err != nil {
		return r.noEntities, err
	}

	entities := []Entity{}
	for rows.Next() {
		row, err := r.scan(rows)
		if err != nil {
			return r.noEntities, err
		}

		entities = append(entities, r.entity(row))
	}

	return entities, nil
}

// Save implements repository.Repository.
func (r *PostgresRepository[Entity, Specification, Row]) Save(ctx context.Context, entity Entity) error {
	row := r.row(entity)

	query, args, err := squirrel.
		Insert(r.tableName).
		Columns(r.columns...).
		Values(r.values(row)...).
		Suffix(r.upsertSuffix).
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

func (r *PostgresRepository[Entity, Specification, Row]) Size(ctx context.Context, specs ...Specification) (uint32, error) {
	var count uint32
	var err error
	builder := squirrel.
		Select("COUNT(id)").
		From(r.tableName).
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

func New[Entity any, Specification any, Row any](opt Option[Entity, Specification, Row]) repository.Repository[Entity, Specification] {
	r := &PostgresRepository[Entity, Specification, Row]{
		dbm:        opt.DatabaseManager,
		filter:     opt.Filter,
		scan:       opt.Scan,
		entity:     opt.Entity,
		row:        opt.Row,
		values:     opt.Values,
		columns:    opt.Columns,
		tableName:  opt.TableName,
		primaryKey: opt.PrimaryKey,
	}

	r.upsertSuffix = r.makeUpsertSuffix()

	return r
}

func (r *PostgresRepository[Entity, Specification, Row]) query(ctx context.Context, args repository.ListArgs[Specification]) (*sql.Rows, error) {
	builder := squirrel.
		Select(r.columns...).
		From(r.tableName).
		Where(r.filter(args.Filters...))
	builder = r.dbm.Paginate(builder, args.Limit, args.Offset)
	queryStr, queryArgs, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	return r.dbm.Querier(ctx).QueryContext(ctx, queryStr, queryArgs...)
}

func (r *PostgresRepository[Entity, Specification, Row]) makeUpsertSuffix() string {
	parts := make([]string, 0, len(r.columns))
	for _, col := range r.columns {
		parts = append(parts, fmt.Sprintf("%s = excluded.%s", col, col))
	}

	return fmt.Sprintf("ON CONFLICT (%s) DO UPDATE SET %s", r.primaryKey, strings.Join(parts, ", "))
}
