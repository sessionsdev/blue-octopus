package redis

import (
	"crypto/sha256"
	"encoding/hex"
)

type RedisKey interface {
	GetKey() string
}

type GenericKey struct {
	Key string
}

func (k *GenericKey) GetKey() string {
	return k.Key
}

type UserDetailsKey struct {
	Email string
}

func (k *UserDetailsKey) GetKey() string {
	hasher := sha256.New()
	hasher.Write([]byte(k.Email))
	return "user:details:" + hex.EncodeToString(hasher.Sum(nil))
}

type UserSessionKey struct {
	SessionID string
}

func (k *UserSessionKey) GetKey() string {
	return "session:" + k.SessionID
}

type UserSavedGameKey struct {
	Email string
}

func (k *UserSavedGameKey) GetKey() string {
	hasher := sha256.New()
	hasher.Write([]byte(k.Email))
	return "user:game:" + hex.EncodeToString(hasher.Sum(nil))
}
