package api

import (
	"errors"
	"fmt"
	"opskrifter-backend/internal/types"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

type (
	CrudFunc[T types.Identifiable] func(T) (string, error)
	GetFunc[T types.Identifiable]  func(T) (T, error)
)

type QueryOptions struct {
	Filters map[string]any
	Page    int
	PerPage int
	OrderBy string
}

var ErrMissingParentOrChild = errors.New("missing parent or child tag in struct")

func buildInsertQuery(obj any) (string, []any, string) {
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

func buildUpdateQuery(obj any) (string, []any) {
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

func BuildQuery(tableName string, opts QueryOptions) (string, []any) {
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PerPage < 1 {
		opts.PerPage = 10
	}
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}

	offset := (opts.Page - 1) * opts.PerPage
	var args []any
	query := fmt.Sprintf("SELECT * FROM %s", tableName)

	if len(opts.Filters) > 0 {
		var where []string
		i := 1
		for k, v := range opts.Filters {
			where = append(where, fmt.Sprintf("%s = $%d", k, i))
			args = append(args, v)
			i++
		}
		query += " WHERE " + strings.Join(where, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY %s LIMIT $%d OFFSET $%d",
		opts.OrderBy,
		len(args)+1,
		len(args)+2)
	args = append(args, opts.PerPage, offset)

	return query, args
}

func buildQueryOneToManyByType[T types.Identifiable, E types.OneToMany](obj T, elements []E) (string, error) {
	parent_id := obj.GetID()
	relation_table := elements[0].TableName()

	parts := []string{}
	for _, e := range elements {
		parts = append(parts, fmt.Sprintf("(%s, %s)", parent_id, e.GetChildId()))
	}

	query := strings.Join(parts, ", ")

	first := reflect.ValueOf(elements[0]).Elem()
	childType := first.Type()

	parentDB, found_parent := childType.FieldByName("parent")
	childDB, found_child := childType.FieldByName("child")

	if !found_parent || !found_child {
		return "", ErrMissingParentOrChild
	}

	parent_col := parentDB.Tag.Get("db")
	child_col := childDB.Tag.Get("db")

	sql := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES %s", relation_table, parent_col, child_col, query)

	return sql, nil
}
