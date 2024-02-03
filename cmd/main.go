package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/sessionsdev/blue-octopus/internal/handlers"
)

func main() {
	http.HandleFunc("/", handlers.ServeHome)
	http.HandleFunc("/about", handlers.ServeAbout)
	http.HandleFunc("/test", handlers.ServeTestPage)
	http.HandleFunc("/api/game-state", handlers.HandleGameState)
	http.HandleFunc("/api/hello-world", handlers.ServeHelloWorldAPI)
	http.HandleFunc("/api/process-command", handlers.HandleGameCommand)

	// handle static files
	staticPath := filepath.Join(".", "static")
	fs := http.FileServer(http.Dir(staticPath))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server is running at http://localhost:8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
