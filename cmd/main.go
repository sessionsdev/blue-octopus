package main

import (
	"fmt"
	"net/http"
	"path/filepath"
    "github.com/sessionsdev/blue-octopus/internal/templatemanager"
	"github.com/sessionsdev/blue-octopus/internal/handlers"
)

	


func main() {
    tmplManager := templatemanager.NewTemplateManager("./templates")

    http.HandleFunc("/", handlers.ServeHome(tmplManager))
    http.HandleFunc("/api/hello-world", handlers.ServeHelloWorldAPI)

    // Handle static files
    
    staticPath := filepath.Join(".", "static")
    absStaticPath, _ := filepath.Abs(staticPath)
    fmt.Println("Serving static files from:", absStaticPath)

    fs := http.FileServer(http.Dir(staticPath))
    http.Handle("/static/", http.StripPrefix("/static/", fs))


    fmt.Println("Server is running at http://localhost:8090")
    http.ListenAndServe(":8090", nil)
}



