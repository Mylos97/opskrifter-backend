package types

type Identifiable interface {
	GetID() string
	TableName() string
}

type ManyToMany interface {
	GetChildID() string
	TableName() string
}

type OneToMany interface {
	TableName() string
}

type HasManyToMany interface {
	GetManyToMany() [][]ManyToMany
}

type HasOneToMany interface {
	GetOneToMany() [][]OneToMany
}

type IdentifiableWithRelations interface {
	Identifiable
	HasManyToMany
	HasOneToMany
}
