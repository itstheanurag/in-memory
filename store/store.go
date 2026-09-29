package store

type Storer interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
	Len() int
	Keys() []string
	Incr(key string) (int, error)
	Decr(key string) (int, error)
}
