package types

type Identifiable interface {
	GetID() string
	TableName() string
}

type HasManyToMany interface {
	GetManyToMany() [][]ManyToMany
}

type ManyToMany interface {
	GetChildID() string
	TableName() string
}

type IdentifiableWithRelations interface {
	Identifiable
	HasManyToMany
}
