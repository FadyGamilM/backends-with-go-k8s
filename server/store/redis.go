package store

type RedisStore struct {
}

func NewRedisStore(url string) *RedisStore {
	return &RedisStore{}
}

func (rs *RedisStore) Put(k string, v []byte) error {
	return nil
}

func (rs *RedisStore) Get(k string) ([]byte, error) {
	return nil, nil
}

func (rs *RedisStore) Ping() error {
	return nil
}
