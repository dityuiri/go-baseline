package db

//go:generate mockgen -destination=mock/rows.go -package=mock . IRows

import "database/sql"

// Rows is the result of a query. Its cursor starts before the first row
// of the result set. Use Next to advance from row to row.
type Rows struct {
	rows *sql.Rows
}

// IRows defines the functions of this interface.
type IRows interface {
	Close() error

	Next() bool
	Scan(dest ...interface{}) error

	Err() error
}

// Close returns the connection to the connection pool.
func (i *Rows) Close() error {
	return i.rows.Close()
}

// Next prepares the next result row for reading with the Scan method.
func (i *Rows) Next() bool {
	return i.rows.Next()
}

// Scan copies the columns from the matched row into the values pointed at by dest.
func (i *Rows) Scan(dest ...interface{}) error {
	return i.rows.Scan(dest...)
}

// Err returns the error, if any, that was encountered during iteration.
func (i *Rows) Err() error {
	return i.rows.Err()
}
