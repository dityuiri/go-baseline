package db

//go:generate mockgen -destination=mock/transaction.go -package=mock . ITransaction

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

// Transaction contains information about this module.
type Transaction struct {
	db *Database
	tx *sql.Tx
}

// ITransaction defines the functions of this module.
type ITransaction interface {
	IEntity

	Commit() error
	Rollback() error
}

// Commit commits the transaction.
func (i *Transaction) Commit() error {
	if err := i.tx.Commit(); err != nil {
		_ = i.Rollback()

		return err
	}

	return nil
}

// Rollback aborts the transaction.
func (i *Transaction) Rollback() error {
	return i.tx.Rollback()
}

// Execute executes a query without returning any rows.
func (i *Transaction) Execute(query string, args ...interface{}) (IResult, error) {
	return i.tx.ExecContext(i.db.context, query, args...)
}

// Query executes a query that returns rows.
func (i *Transaction) Query(query string, args ...interface{}) (IRows, error) {
	rows, err := i.tx.QueryContext(i.db.context, query, args...)

	return &Rows{rows: rows}, err
}

// QueryRow executes a query that returns rows.
func (i *Transaction) QueryRow(query string, args ...interface{}) IRow {
	row := i.tx.QueryRowContext(i.db.context, query, args...)

	return &Row{row: row}
}

// ExecuteContext executes a query on custom context without returning any rows.
func (i *Transaction) ExecuteContext(ctx context.Context, query string, args ...interface{}) (IResult, error) {
	return i.tx.ExecContext(ctx, query, args...)
}

// QueryContext executes a query on custom context that returns rows.
func (i *Transaction) QueryContext(ctx context.Context, query string, args ...interface{}) (IRows, error) {
	rows, err := i.tx.QueryContext(ctx, query, args...)

	return &Rows{rows: rows}, err
}

// QueryRowContext executes a query on custom context that returns rows.
func (i *Transaction) QueryRowContext(ctx context.Context, query string, args ...interface{}) IRow {
	row := i.tx.QueryRowContext(ctx, query, args...)

	return &Row{row: row}
}

// New will insert the entity into the database
func (i *Transaction) New(entity interface{}) (*uuid.UUID, error) {
	return i.db.New(
		entity,
	)
}

// Add will upsert the entity into the database
func (i *Transaction) Add(entity interface{}, columnExcludes ...string) (*uuid.UUID, error) {
	return i.db.Add(
		entity,
		columnExcludes...,
	)
}

// Update will update the entity in the database
func (i *Transaction) Update(entity interface{}, columnExcludes ...string) (*uuid.UUID, error) {
	return i.db.Update(
		entity,
		columnExcludes...,
	)
}

// Get will select the entity from the database
func (i *Transaction) Get(entity interface{}, columnClauses ...string) (interface{}, error) {
	return i.db.Get(entity, columnClauses...)
}

// GetStatement will select the entity from the database using additional statement
func (i *Transaction) GetStatement(entity interface{}, statement string, columnClauses ...string) (interface{}, error) {
	return i.db.GetStatement(
		entity,
		statement,
		columnClauses...,
	)
}

// GetRows will select the entities from the database
func (i *Transaction) GetRows(entity interface{}, columnClauses ...string) (interface{}, error) {
	return i.db.GetRows(entity, columnClauses...)
}

// GetRowsStatement will select the entities from the database using additional statement
func (i *Transaction) GetRowsStatement(entity interface{}, statement string, columnClauses ...string) (interface{}, error) {
	return i.db.GetRowsStatement(
		entity,
		statement,
		columnClauses...,
	)
}

// Remove will delete the entity from the database
func (i *Transaction) Remove(entity interface{}, columnClauses ...string) (*uuid.UUID, error) {
	return i.db.Remove(
		entity,
		columnClauses...,
	)
}
