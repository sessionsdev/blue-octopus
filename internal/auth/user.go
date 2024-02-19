package auth

import (
	"log"
	"net/mail"

	"github.com/sessionsdev/blue-octopus/internal/redis"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username     string
	PasswordHash string
	Email        string
	Role         string
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) IsAdmin() bool {
	return u.Role == "ADMIN"
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	log.Println("Hashed password: ", string(bytes))
	return string(bytes)
}

func CreateUser(username string, password string, email string) *User {
	if GetUser(username) != nil {
		log.Println("User already exists: ", username)
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
		Username:     username,
		PasswordHash: hashPassword(password),
		Role:         "USER",
	}

	redis.SetObj("user", username, user, 0)
	return user
}

func GetUser(username string) *User {
	log.Println("Looking for user: ", username)

	var user User
	_, err := redis.GetObj("user", username, &user)
	if err != nil {
		log.Println("User not found: ", username)
		return nil
	}

	log.Println("Found user: ", username)
	return &user
}

func AuthenticateUser(username string, password string) bool {
	user := GetUser(username)
	if user == nil {
		log.Println("User not found: ", username)
		return false
	}
	log.Println("Authenticating user: ", username)

	return user.CheckPassword(password)
}

func CreateAdminUser(username string, password string, email string) {
	if GetUser(username) != nil {
		log.Println("Admin user already exists: ", username)
		return
	}
	user := &User{
		Username:     username,
		PasswordHash: hashPassword(password),
		Email:        email,
		Role:         "ADMIN",
	}

	redis.SetObj("user", username, user, 0)
}

func IsUserAdmin(username string) bool {
	log.Println("Checking if user is admin: ", username)
	user := GetUser(username)
	if user == nil {
		return false
	}
	return user.IsAdmin()
}
