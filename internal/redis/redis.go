package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"os"
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
	// Get env variables
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	// populate name to client map
	for name, dbId := range dbnameToDbId {
		dbClient := redis.NewClient(&redis.Options{
			Addr:     redisHost + ":" + redisPort,
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

func SetObj(dbName string, key string, value interface{}, expMinutes int) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(value)
	if err != nil {
		return err
	}

	bytes := buf.Bytes()

	client := dbNameToClient[dbName]
	err = client.Set(
		cxt,
		key,
		bytes,
		time.Duration(expMinutes)*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetObj(dbName string, key string, target interface{}) (interface{}, error) {
	client := dbNameToClient[dbName]
	data, err := client.Get(cxt, key).Bytes()
	if err == redis.Nil {
		log.Println("Key not found: ", key)
		return nil, &NotFoundError{Key: key}
	} else if err != nil {
		log.Println("Error getting: ", err)
		return nil, err
	}

	dec := gob.NewDecoder(bytes.NewReader(data))
	err = dec.Decode(target)
	if err != nil {
		log.Println("Error decoding: ", err)
		return nil, err
	}

	return target, nil
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
	value, err := client.Get(cxt, key).Result()
	if err == redis.Nil {
		return "", &NotFoundError{Key: key}
	} else if err != nil {
		return "", err
	}

	return value, nil
}

func DeleteKey(dbName string, key string) error {
	client := dbNameToClient[dbName]
	err := client.Del(cxt, key).Err()
	if err != nil {
		return err
	}
	return nil
}

// get all keys from a db
func GetAllKeys(dbName string) ([]string, error) {
	client := dbNameToClient[dbName]
	keys, err := client.Keys(cxt, "*").Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}
