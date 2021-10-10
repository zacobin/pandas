package kvstore

type KvStore interface {
	Open() error
	Close() error
	Set(key string, value interface{}) error
	Replace(key string, value interface{}) error
	Get(key string) (interface{}, bool)
	Delete(key string) error
	Keys() (keys []string, err error)
	Clean() error
}

func GetKvStore(path string) KvStore {
	// return default kvstore now
	return simpleStore.Load(path)
}
