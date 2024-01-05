package db

//go:generate mockgen -destination=mock/database.go -package=mock . IDatabase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL
)

// Database contains information about this module.
type Database struct {
	context context.Context
	config  *Configuration

	db *sql.DB
}

// IDatabase defines the functions of this module.
type IDatabase interface {
	IEntity

	Close() error
	Ping() error

	Begin(opts ...*sql.TxOptions) (ITransaction, error)
}

// SQLOpenFn is an alias to override the SQL open function.
type SQLOpenFn func(driverName string, dataSourceName string) (*sql.DB, error)

// NewDatabase returns an instance to the database.
func NewDatabase(ctx context.Context, config *Configuration, args ...SQLOpenFn) (IDatabase, error) {
	driverName := config.Driver

	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s password='%s' dbname=%s search_path=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Name,
		config.Schema,
		config.SSLMode,
	)

	sqlOpen := sql.Open
	if len(args) > 0 {
		sqlOpen = args[0]
	}

	if db, err := sqlOpen(driverName, dataSourceName); err != nil {
		return nil, err
	} else {
		// https://www.alexedwards.net/blog/configuring-sqldb
		maxOpenConns := config.MaxOpenConns
		db.SetMaxOpenConns(maxOpenConns)

		maxIdleConns := config.MaxIdleConns
		db.SetMaxIdleConns(maxIdleConns)

		connMaxLifetime := config.ConnMaxLifetime
		db.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

		connMaxIdleTime := config.ConnMaxIdleTime
		db.SetConnMaxIdleTime(time.Duration(connMaxIdleTime) * time.Second)

		return &Database{
			context: ctx,
			config:  config,

			db: db,
		}, nil
	}
}

// Close closes the database and prevents new queries from starting.
func (i *Database) Close() error {
	return i.db.Close()
}

// Ping verifies a connection to the database is still alive, establishing a connection if necessary.
func (i *Database) Ping() error {
	return i.db.PingContext(i.context)
}

// Begin starts a transaction.
// A transaction must end with a call to Commit or Rollback.
// After a call to Commit or Rollback, all operations on the transaction fail with ErrTxDone.
func (i *Database) Begin(opts ...*sql.TxOptions) (ITransaction, error) {
	var txOpts *sql.TxOptions

	if len(opts) > 0 {
		txOpts = opts[0]
	}

	if tx, err := i.db.BeginTx(i.context, txOpts); err != nil {
		return nil, err
	} else {
		return &Transaction{
			db: i,
			tx: tx,
		}, nil
	}
}
