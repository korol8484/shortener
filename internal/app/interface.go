package app

type Entity struct {
	Url   string
	Alias string
}

type Store interface {
	Add(ent *Entity) error
	Read(alias string) (*Entity, error)
}
