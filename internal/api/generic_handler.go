package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"opskrifter-backend/internal/types"
	"opskrifter-backend/pkg/db"
)

func DeleteByType[T types.Identifiable](obj T) (string, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", obj.TableName())
	sql, err := db.DB.Exec(query, obj.GetID())
	rowsAffected, _ := sql.RowsAffected()
	if rowsAffected != 1 {
		return "", ErrRowsAffectedZero
	}
	return obj.GetID(), err
}

func GetByType[T types.Identifiable](obj T) (T, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", obj.TableName())
	err := db.DB.Get(&obj, query, obj.GetID())
	return obj, err
}

func CreateByType[T types.Identifiable](obj T) (string, error) {
	query, args, id := buildInsertQuery(obj)
	sql, err := db.DB.Exec(query, args...)

	if err != nil {
		fmt.Printf("CreateByType: insert failed: %v\n", err)
		return "", err
	}

	fmt.Printf("id %s", "this is where i die")

	rowsAffected, err := sql.RowsAffected()
	if rowsAffected != 1 {
		return "", ErrRowsAffectedZero
	}
	fmt.Printf("id %s", "okay i escaped")

	return id, err
}

func UpdateByType[T types.Identifiable](obj T) (string, error) {
	query, args := buildUpdateQuery(obj)
	sql, err := db.DB.Exec(query, args...)
	rowsAffected, _ := sql.RowsAffected()
	if rowsAffected != 1 {
		return "", ErrRowsAffectedZero
	}

	return obj.GetID(), err
}

func GetCountByType[T types.Identifiable](obj T) (int, error) {
	count := 0
	println(obj.TableName())
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, obj.TableName())
	err := db.DB.QueryRow(query).Scan(&count)
	return count, err
}

func GetManyByType[T types.Identifiable](opts QueryOptions) ([]T, error) {
	var zero T
	var objs []T

	query, args := BuildQuery(zero.TableName(), opts)

	err := db.DB.Select(&objs, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return objs, nil
}

func CreateOneToManyByType[T types.Identifiable, E types.OneToMany](obj T, elements []E) error {
	if len(elements) == 0 {
		return nil
	}

	sql, err := buildQueryOneToManyByType(obj, elements)
	if err != nil {
		return err
	}

	_, err = db.DB.Exec(sql)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	return nil
}

func HandlerByType[T types.Identifiable](crudFunc CrudFunc[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var obj T

		if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if obj.GetID() == "" {
			http.Error(w, "missing ID field", http.StatusBadRequest)
			return
		}

		id, err := crudFunc(obj)
		if err != nil {
			http.Error(w, "operation failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "operation succeeded on ID: %s\n", id)
	}
}

func GetHandlerByType[T types.Identifiable](getFunc GetFunc[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var obj T

		if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if obj.GetID() == "" {
			http.Error(w, "missing ID field", http.StatusBadRequest)
			return
		}

		result, err := getFunc(obj)
		if err != nil {
			http.Error(w, "operation failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}
