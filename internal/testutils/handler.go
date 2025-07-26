package testutils

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"opskrifter-backend/internal/types"
	"reflect"
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

func EqualByValue[T types.Identifiable](want, got T) error {
	if !reflect.DeepEqual(want, got) {
		return fmt.Errorf("structs are not equal\nWant: %+v\nGot:  %+v", want, got)
	}
	return nil
}
