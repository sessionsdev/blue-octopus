package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/sessionsdev/blue-octopus/internal/redis"
)

type UserSession struct {
	Username  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func SaveSession(username string) (string, error) {
	// validate session
	if username == "" {
		return "", fmt.Errorf("username is required")
	}

	// generate session id
	id, err := getUsernameHash(username)
	if err != nil {
		return "", err
	}

	// save session to redis
	err = redis.SetObj("session", id, username, 24*60)
	if err != nil {
		return "", err
	}

	return id, nil
}

func DeleteSession(sessionId string) error {
	return redis.DeleteKey("session", sessionId)
}

func BuildSessionCookie(sessionId string) *http.Cookie {
	return &http.Cookie{
		Name:    "SESSION_ID",
		Value:   sessionId,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	}
}

func BuildDeleteSessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:   "SESSION_ID",
		Value:  "",
		MaxAge: -1,
	}
}

func ValidateSession(r *http.Request) (string, error) {
	// get session id from cookie
	cookie, err := r.Cookie("SESSION_ID")
	if err != nil {
		return "", fmt.Errorf("session cookie not found")
	}

	// get session from redis
	sessionId := cookie.Value
	username, err := redis.GetValue("session", sessionId)
	if err != nil {
		return "", err
	}

	return username, nil
}

func getUsernameHash(username string) (string, error) {
	// Create a session token by hashing the username and password with a secret key
	hash := sha256.New()
	hash.Write([]byte(username))
	token := hex.EncodeToString(hash.Sum(nil))
	return token, nil
}
