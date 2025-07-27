package api

import (
	"fmt"
	"opskrifter-backend/internal/types"
	"opskrifter-backend/pkg/myDB"
)

func DeleteByType[T types.Identifiable](obj T) (string, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", obj.TableName())
	sqlResult, err := myDB.DB.Exec(query, obj.GetID())
	if err != nil {
		return "", fmt.Errorf("failed to delete: %w", err)
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected != 1 {
		return "", ErrRowsAffectedZero
	}

	return obj.GetID(), nil
}

func GetByType[T types.Identifiable](obj T) (T, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", obj.TableName())
	err := myDB.DB.Get(&obj, query, obj.GetID())
	return obj, err
}

func CreateByType[T types.Identifiable](obj T) (string, error) {
	query, args, id := buildInsertQuery(obj)
	result, err := myDB.DB.Exec(query, args...)
	if err != nil {
		return "", fmt.Errorf("failed to execute insert: %w (query: %q)", err, query)
	}

	if result == nil {
		return "", fmt.Errorf("unexpected nil result from Exec()")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected != 1 {
		return "", fmt.Errorf("%w: expected 1 row affected, got %d", ErrRowsAffectedZero, rowsAffected)
	}

	return id, nil
}

func CreateByTypeWithRelations[T types.IdentifiableWithRelations](obj T) (string, error) {
	id, err := CreateByType(obj)

	if err != nil {
		return "", fmt.Errorf("error creating object")
	}

	if id == "" {
		return "", fmt.Errorf("error generating id")
	}

	relations := obj.GetOneToMany()

	for i := range relations {
		err = CreateOneToManyByType(obj, id, relations[i])
		if err != nil {
			return "", err
		}
	}

	return id, nil
}

func UpdateByType[T types.Identifiable](obj T) (string, error) {
	query, args := buildUpdateQuery(obj)
	sqlResult, err := myDB.DB.Exec(query, args...)

	if err != nil {
		return "", fmt.Errorf("failed to update: %w", err)
	}

	rowsAffected, _ := sqlResult.RowsAffected()
	if rowsAffected != 1 {
		return "", ErrRowsAffectedZero
	}

	return obj.GetID(), err
}

func GetCountByType[T types.Identifiable](obj T) (int, error) {
	count := 0
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, obj.TableName())
	err := myDB.DB.QueryRow(query).Scan(&count)
	return count, err
}

func GetCountByTable(table string) (int, error) {
	count := 0
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, table)
	err := myDB.DB.QueryRow(query).Scan(&count)
	return count, err
}

func GetManyByType[T types.Identifiable](opts QueryOptions) ([]T, error) {
	var zero T
	var objs []T
	if opts.PerPage < 0 {
		return nil, fmt.Errorf("per page cannot be less than 0")
	}

	if opts.Page < 0 {
		return nil, fmt.Errorf("page cannot be less than 0")
	}

	query, args := BuildQuery(zero.TableName(), opts)

	err := myDB.DB.Select(&objs, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return objs, nil
}

func CreateManyByType[T types.Identifiable](elements []T) ([]string, error) {
	var ids []string
	for i := range elements {
		id, err := CreateByType(elements[i])
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func DeleteManyByType[T types.Identifiable](elements []T) error {
	for i := range elements {
		_, err := DeleteByType(elements[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateOneToManyByType[T types.Identifiable, E types.OneToMany](obj T, id string, elements []E) error {
	if len(elements) == 0 {
		return nil
	}

	query, err := buildQueryOneToManyByType(id, elements)
	if err != nil {
		return err
	}

	_, err = myDB.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	return nil
}
