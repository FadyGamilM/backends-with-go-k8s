package store

import (
	"log"
	"time"

	"github.com/go-redis/redis"
)

type RedisStore struct {
	rc *redis.Client
}

func NewRedisStore(url string) *RedisStore {
	return &RedisStore{
		rc: redis.NewClient(&redis.Options{
			Addr:         url,
			Password:     "",
			DB:           0,
			ReadTimeout:  time.Minute,
			WriteTimeout: time.Minute,
			DialTimeout:  time.Minute,
		}),
	}
}

func (rs *RedisStore) Put(k string, v []byte) error {
	statusCmd := rs.rc.Set(k, v, time.Hour*24)
	if statusCmd.Err() != nil {
		log.Printf("error : %v setting k %v with v %v\n", statusCmd.Err(), k, v)
		return statusCmd.Err()
	}
	return nil
}

func (rs *RedisStore) Get(k string) []byte {
	v := rs.rc.Get(k)
	log.Println("==============> ", v.Val())
	return []byte(v.Val())
}

func (rs *RedisStore) Ping() error {
	statusCmd := rs.rc.Ping()
	return statusCmd.Err()
}
