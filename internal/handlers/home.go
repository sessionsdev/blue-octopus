package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func ServeHome(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles(
        filepath.Join("templates", "base.html"),
        filepath.Join("templates", "home.html"),
    )

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
