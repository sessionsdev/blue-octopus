package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sessionsdev/blue-octopus/internal/auth"
	"github.com/sessionsdev/blue-octopus/internal/redis"
	"github.com/sessionsdev/blue-octopus/internal/router"
)

func main() {
	staticPath := filepath.Join(".", "static")
	router.Init(staticPath)
	redis.Init()

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	adminEmail := os.Getenv("ADMIN_EMAIL")
	auth.CreateAdminUser(context.TODO(), adminPassword, adminEmail)

	fmt.Println("Server is running at http://localhost:8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
