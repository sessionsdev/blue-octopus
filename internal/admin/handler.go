package admin

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/sessionsdev/blue-octopus/internal/auth"
	"github.com/sessionsdev/blue-octopus/internal/redis"
)

type AdminData struct {
	Users []AdminUserData
}

type AdminUserData struct {
	Email     string
	Role      string
	EmailHash string
}

func BuildFromAuthUser(user auth.User) AdminUserData {
	return AdminUserData{
		Email:     user.Email,
		Role:      user.Role,
		EmailHash: user.HashEmail(),
	}
}

func CheckIfUserContextIsAdmin(ctx context.Context) bool {
	user := ctx.Value("user")
	if user == nil {
		return false
	}

	userDetails := user.(*auth.User)
	return userDetails.Role == "ADMIN"
}

func ServeAdminPage(w http.ResponseWriter, r *http.Request) {
	if !CheckIfUserContextIsAdmin(r.Context()) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	data := AdminData{Users: getUsers(r.Context())}

	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/admin.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getUsers(ctx context.Context) []AdminUserData {
	userkeys, err := redis.GetAllUserDetailKeys(ctx)
	if err != nil {
		log.Println("Failed to get user keys from redis: ", err)
		return nil
	}

	users := make([]AdminUserData, 0)
	for _, keyString := range userkeys {
		key := &redis.GenericKey{Key: keyString}
		var user auth.User
		_, err := redis.GetObj(ctx, key, &user)
		if err != nil {
			log.Println("Failed to get user from redis: ", err)
			continue
		}
		users = append(users, BuildFromAuthUser(user))
	}

	return users
}

func HandleDeleteUserAction(w http.ResponseWriter, r *http.Request) {
	// delete request only
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	if !CheckIfUserContextIsAdmin(r.Context()) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// get id from the path
	userId := r.URL.Query().Get("id")
	if userId == "" {
		http.Error(w, "Missing user id", http.StatusBadRequest)
		return
	}

	// delete user
	key := &redis.GenericKey{Key: "user:details:" + userId}
	err := redis.DeleteKey(r.Context(), key)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	// if hx request then redirect to admin page
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Add("HX-Redirect", "/admin")

	} else {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

func HandleCreateUserForm(w http.ResponseWriter, r *http.Request) {
	// post request only
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	if !CheckIfUserContextIsAdmin(r.Context()) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// get form values
	email := r.FormValue("email")
	password := r.FormValue("password")

	// check if email and password are not empty
	if email == "" || password == "" {
		http.Error(w, "Missing email or password", http.StatusBadRequest)
		return
	}

	// check if user already exists
	if auth.GetUserByEmail(r.Context(), email) != nil {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	// create user
	newUser := auth.CreateUser(r.Context(), email, password)
	if newUser == nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
