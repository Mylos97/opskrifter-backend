package testutils

import (
	"math/rand"
	"opskrifter-backend/internal/types"
	"reflect"
	"time"
)

type TestDataGenerator[T types.Identifiable] struct {
	rng *rand.Rand
}

func NewTestDataGenerator[T types.Identifiable]() *TestDataGenerator[T] {
	src := rand.NewSource(time.Now().UnixNano())
	return &TestDataGenerator[T]{
		rng: rand.New(src),
	}
}

func (g *TestDataGenerator[T]) Generate() T {
	var item T
	t := reflect.TypeOf(item)
	itemValue := reflect.New(t).Elem()

	numFields := itemValue.NumField()
	for i := range numFields {
		field := itemValue.Field(i)

		if !field.CanSet() {
			continue
		}

		if t.Field(i).Name == "ID" {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(generateRandomString(64))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			field.SetInt(g.rng.Int63n(1000))
		case reflect.Float32, reflect.Float64:
			field.SetFloat(g.rng.Float64() * 100)
		case reflect.Bool:
			field.SetBool(g.rng.Intn(2) == 1)
		case reflect.Array:

		}

	}

	item = itemValue.Interface().(T)
	return item
}

func (g *TestDataGenerator[T]) GenerateMany(count int) []T {
	var items []T
	for range count {
		items = append(items, g.Generate())
	}
	return items
}

var (
	idCharset = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = idCharset[rand.Intn(len(idCharset))]
	}
	return string(b)
}
