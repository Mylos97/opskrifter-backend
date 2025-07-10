package types

type Identifiable interface {
	GetID() string
	TableName() string
}
