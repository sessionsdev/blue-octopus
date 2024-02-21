package router

import (
	"log"
	"net/http"

	"github.com/sessionsdev/blue-octopus/internal/admin"
	"github.com/sessionsdev/blue-octopus/internal/auth"
	"github.com/sessionsdev/blue-octopus/internal/game"
	"github.com/sessionsdev/blue-octopus/internal/profile"
)

func Init(staticPath string) {
	// Initialize the handlers
	initializeStaticRoutes(staticPath)
	initializeAdminRoutes()
	initializeApiRoutes()
	initializeAuthRoutes()
	initializeWebRoutes()
	initializeGameRoutes()
}

// intialize the admin routes
func initializeAdminRoutes() {
	http.Handle("/admin", auth.AdminAuthMiddleware(RequestLoggerMiddleware(http.HandlerFunc(admin.ServeAdminPage))))
	http.Handle("/admin/create-user", auth.AdminAuthMiddleware(RequestLoggerMiddleware(http.HandlerFunc(admin.HandleCreateUserForm))))
	http.Handle("/admin/delete-user", auth.AdminAuthMiddleware(RequestLoggerMiddleware(http.HandlerFunc(admin.HandleDeleteUserAction))))
}

// intialize the ai adventure game routes
func initializeGameRoutes() {
	http.Handle("/game", RequestLoggerMiddleware(http.HandlerFunc(game.ServeGamePage)))
	http.Handle("/game/process-command", auth.AuthMiddleware(RequestLoggerMiddleware(http.HandlerFunc(game.HandleGameCommand))))
	http.Handle("/game/game-state", auth.AuthMiddleware(RequestLoggerMiddleware(http.HandlerFunc(game.HandleGameState))))
	http.Handle("/game/stats-display", RequestLoggerMiddleware(http.HandlerFunc(game.ServeGameStats)))
}

// intialize api routes
func initializeApiRoutes() {
	http.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}

// intialize the public routes
func initializeWebRoutes() {
	http.Handle("/", RequestLoggerMiddleware(http.HandlerFunc(profile.ServeHome)))
}

// intialize the auth routes
func initializeAuthRoutes() {
	http.Handle("/login", RequestLoggerMiddleware(http.HandlerFunc(auth.ServeLogin)))
	http.Handle("/logout", RequestLoggerMiddleware(http.HandlerFunc(auth.HandleLogout)))
	http.Handle("/api/authorize", RequestLoggerMiddleware(http.HandlerFunc(auth.HandleLogin)))

}

// intialize the static files
func initializeStaticRoutes(staticPath string) {
	// handle static files
	fs := http.FileServer(http.Dir(staticPath))
	http.Handle("/static/", RequestLoggerMiddleware(http.StripPrefix("/static/", fs)))
}

func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
