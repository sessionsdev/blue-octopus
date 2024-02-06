package handlers

import (
	"net/http"
)

func Init(staticPath string) {
	// Initialize the handlers
	initializeApiRoutes()
	initializeWebRoutes()
	initializeStaticRoutes(staticPath)
}

func initializeApiRoutes() {
	http.HandleFunc("/api/hello-world", ServeHelloWorldAPI)
	http.HandleFunc("/api/process-command", HandleGameCommand)
	http.HandleFunc("/api/game-state", HandleGameState)
	http.HandleFunc("/api/location-tree", HandleWorldTreeVisualization)
	http.HandleFunc("/api/stats-display", ServeGameStats)
}

func initializeWebRoutes() {
	http.HandleFunc("/", ServeHome)
	http.HandleFunc("/about", ServeAbout)
	http.HandleFunc("/game", ServeGamePage)
	http.HandleFunc("/test", ServeTestPage)
}

func initializeStaticRoutes(staticPath string) {
	// handle static files
	fs := http.FileServer(http.Dir(staticPath))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}
