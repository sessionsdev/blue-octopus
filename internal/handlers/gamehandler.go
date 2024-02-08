package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/sessionsdev/blue-octopus/internal/game"
)

type Command struct {
	Command string `json:"command"`
}

type UserPromptWithState struct {
	Prompt               string `json:"prompt"`
	ProposedStateChanges string `json:"proposed_state_changes"`
}

type PreparedStats struct {
	Location  string
	Inventory []string
}

var PreparedStatsCache *PreparedStats = &PreparedStats{}

func populatePreparedStatsCache(g *game.Game) {
	PreparedStatsCache.Location = g.World.CurrentLocation.LocationName
	PreparedStatsCache.Inventory = g.Player.Inventory

	if len(PreparedStatsCache.Inventory) == 0 {
		PreparedStatsCache.Inventory = []string{"You're not carrying anything."}
	}
}

func clearPreparedStatsCache() {
	PreparedStatsCache = &PreparedStats{}
}

func getGameIdCookieValue(r *http.Request) string {
	cookie, err := r.Cookie("GameId")
	if err != nil {
		return ""
	}

	return cookie.Value
}

func setGameIdCookie(w http.ResponseWriter, gameId string) {
	twentyFourHours := 24 * time.Hour

	cookie := http.Cookie{
		Name:   "GameId",
		Value:  gameId,
		MaxAge: int(twentyFourHours),
	}

	http.SetCookie(w, &cookie)
}

func clearGameIdCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:   "GameId",
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(w, &cookie)
}

func ServeGamePage(w http.ResponseWriter, r *http.Request) {
	var g *game.Game

	// If game id gameIdFromCookie exists
	gameIdFromCookie := getGameIdCookieValue(r)
	if gameIdFromCookie == "" {
		// If no game id cookie exists, create a new game
		g = game.InitializeNewGame()
		g.SaveGameToRedis()
		setGameIdCookie(w, g.GameId)
	}

	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/header.html",
		"templates/game.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleGameCommand(w http.ResponseWriter, r *http.Request) {
	// Only POST requests are allowed
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	// Decode the request body into a Command struct
	command := r.FormValue("command")
	if command == "" {
		http.Error(w, "Missing command", http.StatusBadRequest)
		return
	}

	g, _ := LoadGameOrInitializeNew(r)
	result, err := g.ProcessGameCommand(command)
	if err != nil {
		executeTemplate(w, "templates/error-update.html", "game-update", result)
	} else {
		w.Header().Set("HX-Trigger-After-Settle", "stats-update")
		executeTemplate(w, "templates/game-update.html", "game-update", struct {
			PlayerCommand      string
			GameMasterResponse string
		}{
			PlayerCommand:      command,
			GameMasterResponse: result,
		})

		setGameIdCookie(w, g.GameId)
		populatePreparedStatsCache(g)
		g.SaveGameToRedis()
	}
}

func ServeGameStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	executeTemplate(w, "templates/stats-panel.html", "stats-panel", PreparedStatsCache)
	clearPreparedStatsCache()
}

func HandleGameState(w http.ResponseWriter, r *http.Request) {
	// For GET requests, return the current full game state
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
	}
	g, _ := LoadGameOrInitializeNew(r)

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

func LoadGameOrInitializeNew(r *http.Request) (*game.Game, error) {
	// get game id cookie
	cookie, err := r.Cookie("GameId")
	if err != nil {
		return game.InitializeNewGame(), nil
	}

	g, err := game.LoadGameFromRedis(cookie.Value)
	if err != nil || g == nil {
		return game.InitializeNewGame(), nil
	}

	return g, nil
}
