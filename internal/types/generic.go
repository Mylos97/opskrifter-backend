package types

type Identifiable interface {
	GetID() string
	TableName() string
	GetOneToMany() [][]OneToMany
}

type OneToMany interface {
	GetChildId() string
	TableName() string
}
