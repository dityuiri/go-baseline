package db

import "context"

//go:generate mockgen -destination=mock/query.go -package=mock . IQuery

type IQuery interface {
	Execute(query string, args ...interface{}) (IResult, error)

	Query(query string, args ...interface{}) (IRows, error)
	QueryRow(query string, args ...interface{}) IRow

	ExecuteContext(ctx context.Context, query string, args ...interface{}) (IResult, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (IRows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) IRow
}

// Execute executes a query without returning any rows.
func (i *Database) Execute(query string, args ...interface{}) (IResult, error) {
	return i.db.ExecContext(i.context, query, args...)
}

// Query executes a query that returns rows.
func (i *Database) Query(query string, args ...interface{}) (IRows, error) {
	rows, err := i.db.QueryContext(i.context, query, args...)

	return &Rows{rows: rows}, err
}

// QueryRow executes a query that returns rows.
func (i *Database) QueryRow(query string, args ...interface{}) IRow {
	row := i.db.QueryRowContext(i.context, query, args...)

	return &Row{row: row}
}

// ExecuteContext executes a query on custom context without returning any rows.
func (i *Database) ExecuteContext(ctx context.Context, query string, args ...interface{}) (IResult, error) {
	return i.db.ExecContext(ctx, query, args...)
}

// QueryContext executes a query on custom context that returns rows.
func (i *Database) QueryContext(ctx context.Context, query string, args ...interface{}) (IRows, error) {
	rows, err := i.db.QueryContext(ctx, query, args...)

	return &Rows{rows: rows}, err
}

// QueryRowContext executes a query that returns rows.
func (i *Database) QueryRowContext(ctx context.Context, query string, args ...interface{}) IRow {
	row := i.db.QueryRowContext(ctx, query, args...)

	return &Row{row: row}
}
