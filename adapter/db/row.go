package db

//go:generate mockgen -destination=mock/row.go -package=mock . IRow

import "database/sql"

// Row is the result of calling QueryRow to select a single row.
type Row struct {
	row *sql.Row
}

// IRow defines the functions of this interface.
type IRow interface {
	Scan(dest ...interface{}) error

	Err() error
}

// Scan copies the columns from the matched row into the values pointed at by dest.
func (i *Row) Scan(dest ...interface{}) error {
	return i.row.Scan(dest...)
}

// Err provides a way for wrapping packages to check for query errors without calling Scan.
func (i *Row) Err() error {
	return i.row.Err()
}
