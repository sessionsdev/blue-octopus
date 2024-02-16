package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/sessionsdev/blue-octopus/internal/handlers"
	"github.com/sessionsdev/blue-octopus/internal/redis"
)

func main() {
	staticPath := filepath.Join(".", "static")
	handlers.Init(staticPath)
	redis.Init()

	fmt.Println("Server is running at http://localhost:8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
