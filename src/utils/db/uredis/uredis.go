package uredis

import (
	"context"
	"gin-elastic-percolator/src/config"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisConnection struct {
	addr      string
	password  string
	db        int
	client    *redis.Client
	autoClose bool
}

func NewRedisConnection(addr, password string, db int) *RedisConnection {
	return &RedisConnection{
		addr:      addr,
		password:  password,
		db:        db,
		autoClose: true,
	}
}

func NewRedisConnection_ENV() *RedisConnection {
	db, _ := strconv.Atoi(os.Getenv(config.ENV_RDS_DB))
	return NewRedisConnection(
		os.Getenv(config.ENV_RDS_ADDR),
		os.Getenv(config.ENV_RDS_PASS),
		db,
	)
}

func (o *RedisConnection) connect() {
	o.client = redis.NewClient(
		&redis.Options{
			Addr:     o.addr,
			Password: o.password,
			DB:       o.db,
		},
	)

}

func (o *RedisConnection) Close() error {
	if !o.autoClose {
		log.Println("closing redis connection")
	}
	return o.client.Close()
}

// Don't forget to close connection with above function ↑ if u're using this method
func (o *RedisConnection) DisableAutoClose() {
	o.connect()
	o.autoClose = false
}

func (o *RedisConnection) Get(key string) (string, error) {
	if o.autoClose {
		o.connect()
		defer o.Close()
	}

	return o.client.Get(context.Background(), key).Result()
}

func (o *RedisConnection) Set(key, value string, expiration time.Duration) error {
	if o.autoClose {
		o.connect()
		defer o.Close()
	}

	return o.client.Set(context.Background(), key, value, expiration).Err()
}

func (o *RedisConnection) IsKeyExist(keys string) (isExist bool) {
	if o.autoClose {
		o.connect()
		defer o.Close()
	}
	keysList, err := o.client.Keys(context.Background(), keys).Result()
	if err != nil {
		return false
	}

	if len(keysList) > 0 {
		return true
	}
	return false
}
