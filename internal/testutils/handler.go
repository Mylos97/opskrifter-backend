package testutils

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"opskrifter-backend/internal/types"
	"reflect"
	"testing"
)

func NewJSONPostRequest(data []byte) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	return req, rec
}

func AssertCountByType[T types.Identifiable](expected int, getFunc func(T) (int, error)) error {
	var obj T

	count, err := getFunc(obj)
	if err != nil {
		return fmt.Errorf("failed to get count: %w", err)
	}

	if count != expected {
		return fmt.Errorf("expected count to be %d, got %d", expected, count)
	}

	return nil
}

func AssertCountByTable(expected int, tableName string, getFunc func(string) (int, error)) error {
	count, err := getFunc(tableName)
	if err != nil {
		return fmt.Errorf("failed to get count: %w", err)
	}

	if count != expected {
		return fmt.Errorf("expected count to be %d, got %d", expected, count)
	}

	return nil
}

func EqualByValue[T types.Identifiable](want, got T) error {
	if !reflect.DeepEqual(want, got) {
		return fmt.Errorf("structs are not equal\nWant: %+v\nGot:  %+v", want, got)
	}
	return nil
}

func AssertSortedBy[T any](t *testing.T, list []T, less func(a, b T) bool) {
	for i := 1; i < len(list); i++ {
		if !less(list[i-1], list[i]) {
			t.Fatalf("list not sorted at index %d: %v > %v", i, list[i-1], list[i])
		}
	}
}
