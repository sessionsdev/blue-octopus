package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/mail"

	"github.com/sessionsdev/blue-octopus/internal/redis"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	PasswordHash string
	Email        string
	Role         string
}

func (u *User) HashEmail() string {
	return getStringHash(u.Email)
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) IsAdmin() bool {
	return u.Role == "ADMIN"
}

func encryptPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	log.Println("Hashed password: ", string(bytes))
	return string(bytes)
}

func CreateUser(ctx context.Context, email string, password string) *User {
	if GetUserByEmail(ctx, email) != nil {
		log.Println("User already exists: ", email)
		return nil
	}

	if email != "" {
		// validate email
		_, err := mail.ParseAddress(email)
		if err != nil {
			log.Println("Invalid email address: ", email)
			return nil
		}
	}

	var user *User
	user = &User{
		Email:        email,
		PasswordHash: encryptPassword(password),
		Role:         "USER",
	}

	key := &redis.UserDetailsKey{Email: email}
	redis.SetObj(ctx, key, user, 0)
	return user
}

func GetUserByEmail(ctx context.Context, email string) *User {
	key := &redis.UserDetailsKey{Email: email}

	var user User
	_, err := redis.GetObj(ctx, key, &user)
	if err != nil {
		return nil
	}

	return &user
}

func AuthenticateUser(ctx context.Context, email string, password string) bool {
	user := GetUserByEmail(ctx, email)
	if user == nil {
		return false
	}
	return user.CheckPassword(password)
}

func CreateAdminUser(ctx context.Context, password string, email string) {
	if GetUserByEmail(ctx, email) != nil {
		log.Println("Admin user already exists: ", email)
		return
	}
	user := &User{
		PasswordHash: encryptPassword(password),
		Email:        email,
		Role:         "ADMIN",
	}

	key := &redis.UserDetailsKey{Email: email}
	redis.SetObj(ctx, key, user, 0)
}

func IsUserAdmin(ctx context.Context, email string) bool {
	log.Println("Checking if user is admin: ", email)
	user := GetUserByEmail(ctx, email)
	if user == nil {
		return false
	}
	return user.IsAdmin()
}

func getStringHash(s string) string {
	hasher := sha256.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}
