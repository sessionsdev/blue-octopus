package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/sessionsdev/blue-octopus/internal/auth"
)

func ServeLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
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

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	username, password := r.FormValue("username"), r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "Missing username or password", http.StatusBadRequest)
		return
	}

	login := auth.AuthenticateUser(username, password)
	if !login {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := auth.SaveSession(username)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, auth.BuildSessionCookie(token))
	w.Header().Add("HX-Redirect", "/")
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("SESSION_ID")
	if err != nil {
		http.Error(w, "No session to logout", http.StatusBadRequest)
		return
	}

	sessionToken := sessionCookie.Value
	err = auth.DeleteSession(sessionToken)
	if err != nil {
		log.Println("Failed to delete session: ", err)
	}

	http.SetCookie(w, auth.BuildDeleteSessionCookie())
	w.Header().Add("HX-Redirect", "/")
}
