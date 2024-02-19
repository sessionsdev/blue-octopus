package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/sessionsdev/blue-octopus/internal/auth"
)

func Init(staticPath string) {
	// Initialize the handlers
	initializeApiRoutes()
	initializeAuthRoutes()
	initializeWebRoutes()
	initializeStaticRoutes(staticPath)
}

func initializeApiRoutes() {
	http.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.Handle("/api/process-command", AuthMiddleware(RequestLoggerMiddleware(http.HandlerFunc(HandleGameCommand))))
	http.Handle("/api/game-state", AuthMiddleware(RequestLoggerMiddleware(http.HandlerFunc(HandleGameState))))
	http.Handle("/api/stats-display", RequestLoggerMiddleware(http.HandlerFunc(ServeGameStats)))
	http.Handle("/api/authorize", RequestLoggerMiddleware(http.HandlerFunc(HandleLogin)))
}

func initializeWebRoutes() {
	http.Handle("/", RequestLoggerMiddleware(http.HandlerFunc(ServeHome)))
	http.Handle("/game", RequestLoggerMiddleware(http.HandlerFunc(ServeGamePage)))
}

func initializeAuthRoutes() {
	http.Handle("/login", RequestLoggerMiddleware(http.HandlerFunc(ServeLogin)))
	http.Handle("/logout", RequestLoggerMiddleware(http.HandlerFunc(HandleLogout)))

}

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

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		username, err := auth.ValidateSession(r)
		if err != nil {
			w.Header().Add("HX-Redirect", "/login")
			return
		}

		ctx := context.WithValue(r.Context(), "username", username)

		// User is authenticated, proceed with the request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
