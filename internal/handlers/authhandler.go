package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	"net/http"
	"time"

	"github.com/sessionsdev/blue-octopus/internal/redis"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !userIsAuthenticated(r) {
			w.Header().Add("HX-Redirect", "/login")
			return
		}
		// User is authenticated, proceed with the request
		next.ServeHTTP(w, r)
	})
}

func userIsAuthenticated(r *http.Request) bool {
	sessionCookie, err := r.Cookie("SESSIONS_ID")
	if err != nil {
		return false
	}

	sessionToken := sessionCookie.Value
	username, err := redis.GetValue("session", sessionToken)
	if err != nil || username == "" {
		return false
	}

	return true
}

func ServeLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/header.html",
		"templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleAuthorization(w http.ResponseWriter, r *http.Request) {
	username, password := r.FormValue("username"), r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "Missing username or password", http.StatusBadRequest)
		return
	}

	// hash the password
	hash := sha256.New()
	hash.Write([]byte(password))
	password = hex.EncodeToString(hash.Sum(nil))

	userPassword, err := redis.GetValue("user", username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if userPassword != password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := getEncodedToken(username, password)
	if err != nil {
		http.Error(w, "Error encoding token", http.StatusInternalServerError)
		return
	}

	redis.SetValue("session", token, username, 24*60)
	http.SetCookie(w, &http.Cookie{
		Name:    "SESSIONS_ID",
		Value:   token,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	})

	w.Header().Add("HX-Redirect", "/")
}

func getEncodedToken(username, password string) (string, error) {
	// Implement your authentication logic here
	// This is just a placeholder
	hash := sha256.New()
	hash.Write([]byte(username + password))
	token := hex.EncodeToString(hash.Sum(nil))
	return token, nil
}
