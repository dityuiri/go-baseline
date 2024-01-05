package db

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Helper functions

func (i *Database) insertQuery(entity interface{}) (*uuid.UUID, error) {
	statement := `
		INSERT INTO
			%[1]s
			(%[2]s)
		VALUES
			(%[3]s)
		RETURNING
			"id"
		;`

	return i.executeQuery(entity, statement)
}

func (i *Database) upsertQuery(entity interface{}, columnExcludes ...string) (*uuid.UUID, error) {
	statement := `
		INSERT INTO
			%[1]s
			(%[2]s)
		VALUES
			(%[3]s)
		ON CONFLICT (%[4]s)
		DO %[6]s
			%[5]s
		RETURNING
			"id"
		;`

	return i.executeQuery(entity, statement, columnExcludes...)
}

func (i *Database) updateQuery(entity interface{}, columnExcludes ...string) (*uuid.UUID, error) {
	statement := `
		UPDATE
			%[1]s
		SET
			%[5]s
		WHERE
			"id" = $1::uuid
		RETURNING
			"id"
		;`

	return i.executeQuery(entity, statement, columnExcludes...)
}

func (i *Database) selectQuery(entity interface{}, statement string, columnClauses ...string) error {
	statement = fmt.Sprintf(
		`
		SELECT
			%%[4]s
		FROM
			%%[1]s
		WHERE
			%%[3]s
		%s
		;`,

		statement,
	)

	return i.readQuery(entity, statement, columnClauses...)
}

func (i *Database) selectRowsQuery(entity interface{}, statement string, columnClauses ...string) (interface{}, error) {
	statement = fmt.Sprintf(
		`
		SELECT
			%%[4]s
		FROM
			%%[1]s
		WHERE
			%%[3]s
		ORDER BY
			%%[5]s
		%s
		;`,

		statement,
	)

	return i.readRowsQuery(entity, statement, columnClauses...)
}

func (i *Database) deleteQuery(entity interface{}, columnClauses ...string) (*uuid.UUID, error) {
	statement := `
		DELETE FROM
			%[1]s
		WHERE
			%[3]s
		RETURNING
			"id"
		;`

	return i.removeQuery(entity, statement, columnClauses...)
}

func (i *Database) readQuery(entity interface{}, statement string, columnClauses ...string) error {
	query, columns, values := i.buildReadQuery(entity, statement, columnClauses...)

	if err := i.QueryRow(query,
		values...,
	).Scan(columns...); err != nil {
		return err
	}

	return nil
}

func clone(entity interface{}) reflect.Value {
	copy := reflect.New(reflect.TypeOf(entity).Elem())

	value := reflect.ValueOf(entity).Elem()
	for i := 0; i < value.NumField(); i++ {
		copy.Elem().Field(i).Set(value.Field(i))
	}

	return copy
}

func (i *Database) readRowsQuery(entity interface{}, statement string, columnClauses ...string) (interface{}, error) {
	query, columns, values := i.buildReadQuery(entity, statement, columnClauses...)

	entities := reflect.New(reflect.SliceOf(reflect.TypeOf(entity))).Elem()

	if rows, err := i.Query(query,
		values...,
	); err != nil {
		return nil, err
	} else {
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(columns...); err != nil {
				return nil, err
			}

			entities = reflect.Append(entities, clone(entity))
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	return entities.Interface(), nil
}

func (i *Database) executeQuery(entity interface{}, statement string, columnExcludes ...string) (*uuid.UUID, error) {
	query, _, values := i.buildQuery(entity, statement, columnExcludes...)

	var id *uuid.UUID

	if err := i.QueryRow(query,
		values...,
	).Scan(&id); err != nil {
		return nil, err
	}

	return id, nil
}

func (i *Database) removeQuery(entity interface{}, statement string, columnClauses ...string) (*uuid.UUID, error) {
	query, _, values := i.buildReadQuery(entity, statement, columnClauses...)

	var id *uuid.UUID

	if err := i.QueryRow(query,
		values...,
	).Scan(&id); err != nil {
		return nil, err
	}

	return id, nil
}

func (i *Database) buildQuery(entity interface{}, statement string, columnExcludes ...string) (string, []interface{}, []interface{}) {
	var (
		columns []interface{}
		values  []interface{}

		columnNames  []string
		columnValues []string

		columnUpdates []string
	)

	elem := reflect.ValueOf(entity).Elem()
	tableName := ToSnakeCase(elem.Type().Name())

	withImmutables := strings.Contains(statement, "%[2]s") || strings.Contains(statement, "%[3]s")

	for n := 0; n < elem.NumField(); n++ {
		field := elem.Field(n)
		//fmt.Printf("KIND[%v]\tTYPE[%v]\tVALUE[%v]\n", field.Kind(), field.Type(), field.Interface())

		tag := elem.Type().Field(n).Tag.Get("column")

		//ignore non "column" fields
		if tag == "" {
			continue
		}

		columns = append(columns, field.Addr().Interface())

		column, opts := ParseTag(tag)
		value := field.Interface()

		if opts.Contains("readonly") {
			continue
		}

		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				if column == "id" {
					// Generate UUID for primary key.
					value, _ = uuid.NewUUID()
				} else {
					continue
				}
			}
		} else {
			panic(fmt.Sprintf("column %v must be pointer", column))
		}

		var uuid string
		if field.Type().String() == "uuid.UUID" ||
			field.Type().String() == "*uuid.UUID" {
			uuid = "::uuid"
		}

		if !opts.Contains("immutable") || withImmutables && opts.Contains("immutable") {
			values = append(values, value)

			columnNames = append(columnNames, pq.QuoteIdentifier(column))
			columnValues = append(columnValues, fmt.Sprintf("$%d%s", len(columnNames), uuid))

			if !contains(columnExcludes, column) && column != "id" && !opts.Contains("immutable") {
				columnUpdates = append(columnUpdates, fmt.Sprintf("%s = $%d%s", pq.QuoteIdentifier(column), len(columnNames), uuid))

				if opts.Contains("unique") {
					columnExcludes = append(columnExcludes, pq.QuoteIdentifier(column))
				}
			}
		}
	}

	var query string

	if (len(columnNames) > 0 &&
		len(columnValues) > 0) ||
		len(columns) > 0 {
		var doOperation = "NOTHING"
		if len(columnUpdates) > 0 {
			doOperation = "UPDATE SET"
		}

		if len(columnExcludes) == 0 {
			columnExcludes = []string{"id"}
		}

		query = fmt.Sprintf(
			statement,

			pq.QuoteIdentifier(tableName),

			strings.Join(columnNames, ", "),
			strings.Join(columnValues, ", "),

			strings.Join(columnExcludes, ", "),
			strings.Join(columnUpdates, ", "),

			doOperation,
		)

		if i.config.Echo {
			fmt.Println(query)
		}
	}

	return query, columns, values
}

var sortMapping = map[byte]string{
	'+': "ASC",
	'-': "DESC",
}

func (i *Database) buildReadQuery(entity interface{}, statement string, columnClauses ...string) (string, []interface{}, []interface{}) {
	var (
		columns []interface{}
		values  []interface{}

		columnNames  []string
		columnWheres []string
		columnOrders []string
	)

	processColumnOrders := func(temp []string) []string {
		for _, columnClause := range columnClauses {
			if sort, ok := sortMapping[columnClause[0]]; !ok {
				temp = append(temp, columnClause)
			} else {
				columnOrders = append(columnOrders, pq.QuoteIdentifier(columnClause[1:])+" "+sort)
			}
		}

		return temp
	}

	columnClauses = processColumnOrders(columnClauses[:0])

	elem := reflect.ValueOf(entity).Elem()
	tableName := ToSnakeCase(elem.Type().Name())

	for n := 0; n < elem.NumField(); n++ {
		field := elem.Field(n)
		// fmt.Printf("KIND[%v]\tTYPE[%v]\tVALUE[%v]\n", field.Kind(), field.Type(), field.Interface())

		tag := elem.Type().Field(n).Tag.Get("column")

		//ignore non "column" fields
		if tag == "" {
			continue
		}

		columns = append(columns, field.Addr().Interface())

		column, _ := ParseTag(tag)

		if field.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("column %v must be pointer", column))
		}

		var uuid string
		if field.Type().String() == "uuid.UUID" ||
			field.Type().String() == "*uuid.UUID" {
			uuid = "::uuid"
		}

		columnNames = append(columnNames, pq.QuoteIdentifier(column))

		if (len(columnClauses) == 0 || contains(columnClauses, column)) && !field.IsNil() {
			comparator := "="
			if field.Type().String() == "string" ||
				field.Type().String() == "*string" {
				comparator = "ILIKE"
			}

			columnWheres = append(columnWheres, fmt.Sprintf("%s %s $%d%s", pq.QuoteIdentifier(column), comparator, 1+len(values), uuid))

			values = append(values, field.Interface())
		}
	}

	var query string

	if (len(columnNames) > 0 &&
		len(columnWheres) > 0 &&
		len(columnOrders) > 0) ||
		len(columns) > 0 {
		if len(columnWheres) == 0 {
			columnWheres = append(columnWheres, "NULL")
		}
		if len(columnOrders) == 0 {
			columnOrders = append(columnOrders, "id")
		}

		query = fmt.Sprintf(
			statement,

			pq.QuoteIdentifier(tableName),

			strings.Join(columnNames, ", "),
			strings.Join(columnWheres, " AND "),

			strings.Join(columnNames, ", "),

			strings.Join(columnOrders, ", "),
		)

		if i.config.Echo {
			fmt.Println(query)
		}
	}

	return query, columns, values
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
