package game

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/sessionsdev/blue-octopus/internal/auth"
)

func ServeGameStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
	}

	if PreparedStatsCache == nil {
		http.Error(w, "No stats available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	executeTemplate(w, "templates/stats-panel.html", "stats-panel", PreparedStatsCache)
}

func HandleGameState(w http.ResponseWriter, r *http.Request) {
	// For GET requests, return the current full game state
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
	}

	// get user from context
	user := r.Context().Value("user")
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	email := user.(*auth.User).Email

	g, err := LoadGameFromRedis(r.Context(), email)
	if err != nil {
		http.Error(w, "Error loading game from redis", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func executeTemplate(w http.ResponseWriter, templateFile string, templateName string, data interface{}) {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, templateName, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
