package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client
var cxt = context.Background()

type NotFoundError struct {
	Key string
}

func (e *NotFoundError) Error() string {
	return "Key not found: " + e.Key
}

func Init() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	Client = client

	// ping the redis server to check if it's alive
	pong, err := client.Ping(cxt).Result()
	if err != nil {
		panic(err)
	}

	// print the ping response
	println(pong + " - Redis is alive!")
}

func SetGob(key string, value []byte, expMinutes int) error {
	err := Client.Set(
		cxt,
		key,
		value,
		time.Duration(expMinutes)*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetGob(key string) ([]byte, error) {
	bytes, err := Client.Get(cxt, key).Bytes()
	if err == redis.Nil {
		return nil, &NotFoundError{Key: key}
	} else if err != nil {
		return nil, err
	}

	return bytes, nil
}
