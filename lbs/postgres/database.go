package postgres

import (
	"context"
	"database/sql"

	"github.com/cloustone/pandas/pkg/errors"
	"github.com/cloustone/pandas/users"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	errDuplicate  = "unique_violation"
	errFK         = "foreign_key_violation"
	errInvalid    = "invalid_text_representation"
	errTruncation = "string_data_right_truncation"
)

var _ Database = (*database)(nil)

type database struct {
	db *sqlx.DB
}

// Database provides a database interface
type Database interface {
	NamedExecContext(context.Context, string, interface{}) (sql.Result, error)
	QueryRowxContext(context.Context, string, ...interface{}) *sqlx.Row
	NamedQueryContext(context.Context, string, interface{}) (*sqlx.Rows, error)
	GetContext(context.Context, interface{}, string, ...interface{}) error
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

// NewDatabase creates a ThingDatabase instance
func NewDatabase(db *sqlx.DB) Database {
	return &database{
		db: db,
	}
}

func (dm database) NamedExecContext(ctx context.Context, query string, args interface{}) (sql.Result, error) {
	addSpanTags(ctx, query)
	result, err := dm.db.NamedExecContext(ctx, query, args)
	if pqErr, ok := err.(*pq.Error); ok && errDuplicate == pqErr.Code.Name() {
		return result, errors.Wrap(users.ErrConflict, err)
	}
	return result, err
}

func (dm database) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	addSpanTags(ctx, query)
	return dm.db.QueryRowxContext(ctx, query, args...)
}

func (dm database) NamedQueryContext(ctx context.Context, query string, args interface{}) (*sqlx.Rows, error) {
	addSpanTags(ctx, query)
	return dm.db.NamedQueryContext(ctx, query, args)
}

func (dm database) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	addSpanTags(ctx, query)
	return dm.db.GetContext(ctx, dest, query, args...)
}

func (dm database) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.SetTag("span.kind", "client")
		span.SetTag("peer.service", "postgres")
		span.SetTag("db.type", "sql")
	}
	return dm.db.BeginTxx(ctx, opts)
}

func addSpanTags(ctx context.Context, query string) {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.SetTag("sql.statement", query)
		span.SetTag("span.kind", "client")
		span.SetTag("peer.service", "postgres")
		span.SetTag("db.type", "sql")
	}
}
