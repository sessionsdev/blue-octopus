package handlers

import (
	"log"
	"net/http"
)

func Init(staticPath string) {
	// Initialize the handlers
	initializeApiRoutes()
	initializeWebRoutes()
	initializeStaticRoutes(staticPath)
}

func initializeApiRoutes() {
	http.Handle("/api/hello-world", RequestLoggerMiddleware(http.HandlerFunc(ServeHelloWorldAPI)))
	http.Handle("/api/process-command", RequestLoggerMiddleware(http.HandlerFunc(HandleGameCommand)))
	http.Handle("/api/game-state", RequestLoggerMiddleware(http.HandlerFunc(HandleGameState)))
	http.Handle("/api/stats-display", RequestLoggerMiddleware(http.HandlerFunc(ServeGameStats)))
}

func initializeWebRoutes() {
	http.Handle("/", RequestLoggerMiddleware(http.HandlerFunc(ServeHome)))
	http.Handle("/about", RequestLoggerMiddleware(http.HandlerFunc(ServeAbout)))
	http.Handle("/game", RequestLoggerMiddleware(http.HandlerFunc(ServeGamePage)))
	http.Handle("/test", RequestLoggerMiddleware(http.HandlerFunc(ServeTestPage)))
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
