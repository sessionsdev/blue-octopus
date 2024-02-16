package handlers

import (
	"log"
	"net/http"
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

	http.Handle("/api/hello-world", RequestLoggerMiddleware(http.HandlerFunc(ServeHelloWorldAPI)))
	http.Handle("/api/process-command", AuthMiddleware(RequestLoggerMiddleware(http.HandlerFunc(HandleGameCommand))))
	http.Handle("/api/game-state", RequestLoggerMiddleware(http.HandlerFunc(HandleGameState)))
	http.Handle("/api/stats-display", RequestLoggerMiddleware(http.HandlerFunc(ServeGameStats)))
	http.Handle("/api/authorize", RequestLoggerMiddleware(http.HandlerFunc(HandleAuthorization)))
}

func initializeWebRoutes() {
	http.Handle("/", RequestLoggerMiddleware(http.HandlerFunc(ServeHome)))
	http.Handle("/about", RequestLoggerMiddleware(http.HandlerFunc(ServeAbout)))
	http.Handle("/game", RequestLoggerMiddleware(http.HandlerFunc(ServeGamePage)))
	http.Handle("/test", RequestLoggerMiddleware(http.HandlerFunc(ServeTestPage)))
}

func initializeAuthRoutes() {
	http.Handle("/login", RequestLoggerMiddleware(http.HandlerFunc(ServeLogin)))
	// http.Handle("/logout", RequestLoggerMiddleware(http.HandlerFunc(ServeLogout)))

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
