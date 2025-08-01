package api

import (
	"errors"
	"fmt"
	"opskrifter-backend/internal/types"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

type QueryOptions struct {
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
	OrderBy string `json:"order_by"`
}

var validOrderBys = map[string]bool{
	"id":         true,
	"name":       true,
	"created_at": true,
	"likes":      true,
	"minutes":    true,
}

var ErrMissingParentOrChild = errors.New("missing parent or child tag in struct")
var ErrRowsAffectedZero = errors.New("expected affected rows to be 1 got 0")
var ErrExecutingQuery = errors.New("error executing query")
var ErrNotValidOrderBy = errors.New("this order by does not exist")
var ErrNoColumnNamesFound = errors.New("no column names found")
var ErrNoIdForType = errors.New("no id for type")

func BuildInsertQuery(obj any) (string, []any, string) {
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)
	columns := []string{}
	placeholders := []string{}
	values := []any{}
	id := uuid.New().String()

	for i := range v.NumField() {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			continue
		}
		val := v.Field(i).Interface()

		if dbTag == "id" {
			val = id
		}

		columns = append(columns, dbTag)
		placeholders = append(placeholders, "?")
		values = append(values, val)
	}

	table := obj.(types.Identifiable).TableName()
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	return query, values, id
}

func BuildUpdateQuery(obj any) (string, []any) {
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	assignments := []string{}
	values := []any{}

	var idValue any
	var idColumn string

	for i := range v.NumField() {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			continue
		}

		val := v.Field(i).Interface()
		if dbTag == "id" {
			idValue = val
			idColumn = dbTag
			continue
		}

		assignments = append(assignments, fmt.Sprintf("%s = ?", dbTag))
		values = append(values, val)
	}

	values = append(values, idValue)
	table := obj.(types.Identifiable).TableName()
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", table,
		strings.Join(assignments, ", "),
		idColumn,
	)
	return query, values
}

func BuildQuery(tableName string, opts QueryOptions) (string, []any, error) {
	offset := (opts.Page - 1) * opts.PerPage
	var args []any
	query := fmt.Sprintf("SELECT * FROM %s", tableName)

	if opts.OrderBy != "" && !validOrderBys[opts.OrderBy] {
		return "", nil, ErrNotValidOrderBy
	}

	if opts.OrderBy != "" && validOrderBys[opts.OrderBy] {
		query += fmt.Sprintf(" ORDER BY %s", opts.OrderBy)
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, opts.PerPage, offset)

	return query, args, nil
}

func BuildQueryOneToManyByType[E types.ManyToMany](parentID string, elements []E) (string, []any, error) {
	if len(elements) == 0 {
		return "", nil, fmt.Errorf("no elements provided")
	}

	relationTable := elements[0].TableName()
	first := reflect.ValueOf(elements[0])
	elemType := first.Type()

	var columnNames []string
	var placeholders []string
	var args []any

	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" {
			columnNames = append(columnNames, dbTag)
		}
	}

	for _, element := range elements {
		val := reflect.ValueOf(element)
		var rowPlaceholders []string

		for i := 0; i < elemType.NumField(); i++ {
			field := elemType.Field(i)
			dbTag := field.Tag.Get("db")
			if dbTag == "" {
				continue
			}

			if _, isParent := field.Tag.Lookup("parent"); isParent {
				args = append(args, parentID)
			} else {
				args = append(args, val.Field(i).Interface())
			}

			rowPlaceholders = append(rowPlaceholders, "?")
		}

		placeholders = append(placeholders, fmt.Sprintf("(%s)", strings.Join(rowPlaceholders, ", ")))
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s",
		relationTable,
		strings.Join(columnNames, ", "),
		strings.Join(placeholders, ", "),
	)
	return query, args, nil
}

func GetColumnNames[E types.ManyToMany](element E) ([]string, error) {
	first := reflect.ValueOf(element)
	elemType := first.Type()
	var columnNames []string

	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" {
			columnNames = append(columnNames, dbTag)
		}
	}

	if len(columnNames) == 0 {
		return nil, ErrNoColumnNamesFound
	}

	return columnNames, nil
}
