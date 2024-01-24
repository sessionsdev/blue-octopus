package handlers

import (
	"net/http"
    "github.com/sessionsdev/blue-octopus/internal/templatemanager"
)

func ServeHome(tmplManager *templatemanager.TemplateManager) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        err := tmplManager.ExecuteTemplate(w, "base", nil)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    }
}
