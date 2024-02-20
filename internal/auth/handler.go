package auth

import (
	"context"
	"html/template"
	"log"
	"net/http"
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
	email, password := r.FormValue("username"), r.FormValue("password")
	if email == "" || password == "" {
		http.Error(w, "Missing username or password", http.StatusBadRequest)
		return
	}

	login := AuthenticateUser(email, password)
	if !login {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := SaveSession(email)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, BuildSessionCookie(token))
	w.Header().Add("HX-Redirect", "/")
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("SESSION_ID")
	if err != nil {
		http.Error(w, "No session to logout", http.StatusBadRequest)
		return
	}

	sessionToken := sessionCookie.Value
	err = DeleteSession(sessionToken)
	if err != nil {
		log.Println("Failed to delete session: ", err)
	}

	http.SetCookie(w, BuildDeleteSessionCookie())
	w.Header().Add("HX-Redirect", "/")
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("AuthMiddleware")

		email, err := ValidateSession(r)
		if err != nil {
			handleLoginRedirect(w, r)
			return
		}

		// set user object in context
		user := GetUserByEmail(email)
		if user == nil {
			handleLoginRedirect(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)

		// User is authenticated, proceed with the request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func handleLoginRedirect(w http.ResponseWriter, r *http.Request) {
	// if htmx request, set header
	hxRequest := r.Header.Get("Hx-Request")
	if hxRequest == "true" {
		w.Header().Add("HX-Redirect", "/login")
		w.WriteHeader(http.StatusSeeOther)
		return
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}
