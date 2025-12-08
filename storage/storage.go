package storage

type Data map[string]any

type Storage interface {
	Load(string) (Data, error)
}
