package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type NotFoundError struct {
	Key string
}

func (e *NotFoundError) Error() string {
	return "Key not found: " + e.Key
}

func SetObj(ctx context.Context, key RedisKey, value interface{}, expMinutes int) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(value)
	if err != nil {
		return err
	}

	bytes := buf.Bytes()

	err = Client.Set(
		ctx,
		key.GetKey(),
		bytes,
		time.Duration(expMinutes)*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetObj(ctx context.Context, key RedisKey, target interface{}) (interface{}, error) {
	data, err := Client.Get(ctx, key.GetKey()).Bytes()
	if err == redis.Nil {
		log.Println("Key not found: ", key)
		return nil, &NotFoundError{Key: key.GetKey()}
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

func SetValue(ctx context.Context, key RedisKey, value string, expMinutes int) error {
	err := Client.Set(
		ctx,
		key.GetKey(),
		value,
		time.Duration(expMinutes)*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetValue(ctx context.Context, key RedisKey) (string, error) {
	value, err := Client.Get(ctx, key.GetKey()).Result()
	if err == redis.Nil {
		return "", &NotFoundError{Key: key.GetKey()}
	} else if err != nil {
		return "", err
	}

	return value, nil
}

func DeleteKey(ctx context.Context, key RedisKey) error {
	err := Client.Del(ctx, key.GetKey()).Err()
	if err != nil {
		return err
	}
	return nil
}

// get all keys from a db
func GetAllUserDetailKeys(ctx context.Context) ([]string, error) {
	// TODO - this is not efficient.  Should use scan or a set
	keys, err := Client.Keys(ctx, "user:details:*").Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}
