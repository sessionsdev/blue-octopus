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
	redis.SetValue("user", "admin", "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8", 9999999)

	fmt.Println("Server is running at http://localhost:8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
