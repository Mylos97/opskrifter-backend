package types

type Identifiable interface {
	GetID() string
	TableName() string
}

type HasOneToMany interface {
	GetOneToMany() [][]OneToMany
}

type OneToMany interface {
	GetChildId() string
	TableName() string
}
