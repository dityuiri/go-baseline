package db

//go:generate mockgen -destination=mock/entity.go -package=mock . IEntity

import "github.com/google/uuid"

type IEntity interface {
	IQuery

	New(entity interface{}) (*uuid.UUID, error)
	Add(entity interface{}, columnExcludes ...string) (*uuid.UUID, error)

	Update(entity interface{}, columnExcludes ...string) (*uuid.UUID, error)

	Get(entity interface{}, columnClauses ...string) (interface{}, error)
	GetStatement(entity interface{}, statement string, columnClauses ...string) (interface{}, error)

	GetRows(entity interface{}, columnClauses ...string) (interface{}, error)
	GetRowsStatement(entity interface{}, statement string, columnClauses ...string) (interface{}, error)

	Remove(entity interface{}, columnClauses ...string) (*uuid.UUID, error)
}

// New will insert the entity into the database
func (i *Database) New(entity interface{}) (*uuid.UUID, error) {
	return i.insertQuery(
		entity,
	)
}

// Add will upsert the entity into the database
func (i *Database) Add(entity interface{}, columnExcludes ...string) (*uuid.UUID, error) {
	return i.upsertQuery(
		entity,
		columnExcludes...,
	)
}

// Update will update the entity in the database
func (i *Database) Update(entity interface{}, columnExcludes ...string) (*uuid.UUID, error) {
	return i.updateQuery(
		entity,
		columnExcludes...,
	)
}

// Get will select the entity from the database
func (i *Database) Get(entity interface{}, columnClauses ...string) (interface{}, error) {
	var statement string

	return i.GetStatement(entity, statement, columnClauses...)
}

// GetStatement will select the entity from the database using additional statement
func (i *Database) GetStatement(entity interface{}, statement string, columnClauses ...string) (interface{}, error) {
	if err := i.selectQuery(
		entity,
		statement,
		columnClauses...,
	); err != nil {
		return nil, err
	}

	return entity, nil
}

// GetRows will select the entities from the database
func (i *Database) GetRows(entity interface{}, columnClauses ...string) (interface{}, error) {
	var statement string

	return i.GetRowsStatement(entity, statement, columnClauses...)
}

// GetRowsStatement will select the entities from the database using additional statement
func (i *Database) GetRowsStatement(entity interface{}, statement string, columnClauses ...string) (interface{}, error) {
	if entities, err := i.selectRowsQuery(
		entity,
		statement,
		columnClauses...,
	); err != nil {
		return nil, err
	} else {
		return entities, nil
	}
}

// Remove will delete the entity from the database
func (i *Database) Remove(entity interface{}, columnClauses ...string) (*uuid.UUID, error) {
	return i.deleteQuery(
		entity,
		columnClauses...,
	)
}
