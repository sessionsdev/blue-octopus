package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
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

func SaveSession(ctx context.Context, email string) (string, error) {
	// validate session
	if email == "" {
		return "", fmt.Errorf("username is required")
	}

	// generate session id
	id, err := generateSessionID(32)
	if err != nil {
		return "", err
	}

	// save session to redis
	key := &redis.UserSessionKey{SessionID: id}
	err = redis.SetValue(ctx, key, email, 24*60)
	if err != nil {
		return "", err
	}

	return id, nil
}

func DeleteSession(ctx context.Context, sessionId string) error {
	key := &redis.UserSessionKey{SessionID: sessionId}
	return redis.DeleteKey(ctx, key)
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
	key := &redis.UserSessionKey{SessionID: sessionId}
	username, err := redis.GetValue(r.Context(), key)
	if err != nil {
		return "", err
	}

	// print username
	fmt.Println("Validate sessions - Username: ", username)

	return username, nil
}

func generateSessionID(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
