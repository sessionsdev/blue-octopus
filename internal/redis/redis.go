package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var dbnameToDbId = map[string]int{
	"game":    0,
	"user":    1,
	"session": 2,
}

var dbNameToClient = map[string]*redis.Client{}

var cxt = context.Background()

type NotFoundError struct {
	Key string
}

func (e *NotFoundError) Error() string {
	return "Key not found: " + e.Key
}

func Init() {
	// populate name to client map
	for name, dbId := range dbnameToDbId {
		dbClient := redis.NewClient(&redis.Options{
			Addr:     "red-cn7earo21fec73fld5gg:6379",
			Password: "",
			DB:       dbId,
		})

		pong, err := dbClient.Ping(cxt).Result()
		if err != nil {
			panic(err)
		}

		println(pong + " DB: " + name + " is alive!")
		dbNameToClient[name] = dbClient
	}
}

func SetGob(key string, value []byte, expMinutes int) error {
	client := dbNameToClient["game"]
	err := client.Set(
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
	client := dbNameToClient["game"]
	bytes, err := client.Get(cxt, key).Bytes()
	if err == redis.Nil {
		return nil, &NotFoundError{Key: key}
	} else if err != nil {
		return nil, err
	}

	return bytes, nil
}

func SetValue(dbName string, key string, value string, expMinutes int) error {
	client := dbNameToClient[dbName]
	err := client.Set(
		cxt,
		key,
		value,
		time.Duration(expMinutes)*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetValue(dbName string, key string) (string, error) {
	client := dbNameToClient[dbName]
	token, err := client.Get(cxt, key).Result()
	if err == redis.Nil {
		return "", &NotFoundError{Key: key}
	} else if err != nil {
		return "", err
	}

	return token, nil
}
