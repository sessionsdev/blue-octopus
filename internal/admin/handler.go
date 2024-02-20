package admin

import (
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

func ServeAdminPage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userDetails := user.(*auth.User)
	if userDetails.Role != "ADMIN" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	data := AdminData{Users: getUsers()}

	log.Println("Serving admin page")
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

func getUsers() []AdminUserData {
	usernames, err := redis.GetAllKeys("user")
	if err != nil {
		log.Println("Failed to get user keys from redis: ", err)
		return nil
	}

	users := make([]AdminUserData, 0)
	for _, username := range usernames {
		var user auth.User
		_, err := redis.GetObj("user", username, &user)
		if err != nil {
			log.Println("Failed to get user from redis: ", err)
			continue
		}
		users = append(users, BuildFromAuthUser(user))
	}

	return users
}

func HandleDeleteUser(w http.ResponseWriter, r *http.Request) {

}
