package store

type Storer interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
	Len() int
	Keys() []string
}
