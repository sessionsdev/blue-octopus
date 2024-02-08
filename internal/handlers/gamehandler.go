package handlers

import (
	"encoding/json"
	"html/template"
	"log"
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
	Location           string
	PreviousLocation   string
	PotentialLocations []string
	Inventory          []string
	Enemies            []string
	InteractiveItems   []string
}

var PreparedStatsCache *PreparedStats = &PreparedStats{}

func populatePreparedStatsCache(g *game.Game) {
	PreparedStatsCache.Location = g.World.CurrentLocation.LocationName
	PreparedStatsCache.PreviousLocation = g.World.CurrentLocation.PreviousLocation

	PreparedStatsCache.PotentialLocations = make([]string, 0, len(g.World.CurrentLocation.PotentialLocations))
	for k := range g.World.CurrentLocation.PotentialLocations {
		PreparedStatsCache.PotentialLocations = append(PreparedStatsCache.PotentialLocations, k)
	}

	PreparedStatsCache.Enemies = make([]string, 0, len(g.World.CurrentLocation.Enemies))
	for k := range g.World.CurrentLocation.Enemies {
		PreparedStatsCache.Enemies = append(PreparedStatsCache.Enemies, k)
	}

	PreparedStatsCache.InteractiveItems = make([]string, 0, len(g.World.CurrentLocation.InteractiveItems))
	for k := range g.World.CurrentLocation.InteractiveItems {
		PreparedStatsCache.InteractiveItems = append(PreparedStatsCache.InteractiveItems, k)
	}

	PreparedStatsCache.Inventory = make([]string, 0, len(g.Player.Inventory))
	for k := range g.Player.Inventory {
		PreparedStatsCache.Inventory = append(PreparedStatsCache.Inventory, k)
	}
}

func clearPreparedStatsCache() {
	PreparedStatsCache = &PreparedStats{}
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

	// Decode the request body into a Command struct
	command := r.FormValue("command")
	if command == "" {
		http.Error(w, "Missing command", http.StatusBadRequest)
		return
	}

	var gameId string
	gameIdCookie, err := r.Cookie("GameId")
	if err != nil {
		gameId = ""
	} else {
		gameId = gameIdCookie.Value
	}

	resultMsg, game, err := game.ProcessGameCommand(command, gameId)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		executeTemplate(w, "templates/error-update.html", "game-update", resultMsg)
	} else if game == nil {
		w.Header().Set("Content-Type", "text/html")
		executeTemplate(w, "templates/error-update.html", "game-update", resultMsg)
	} else {
		game.SaveGameToRedis()
		populatePreparedStatsCache(game)

		w.Header().Set("HX-Trigger-After-Settle", "stats-update")
		w.Header().Set("Content-Type", "text/html")

		if gameId != game.GameId {
			setGameIdCookie(w, game.GameId)
		}

		executeTemplate(w, "templates/game-update.html", "game-update", struct {
			PlayerCommand      string
			GameMasterResponse string
		}{
			PlayerCommand:      command,
			GameMasterResponse: resultMsg,
		})
	}
}

func ServeGameStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	executeTemplate(w, "templates/stats-panel.html", "stats-panel", PreparedStatsCache)
}

func HandleGameState(w http.ResponseWriter, r *http.Request) {
	// For GET requests, return the current full game state
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
	}

	var gameId string
	gameIdCookie, err := r.Cookie("GameId")
	if err != nil {
		log.Printf("Error getting game id from cookie: %s", err)
		http.Error(w, "No game id found", http.StatusBadRequest)
		return
	} else {
		gameId = gameIdCookie.Value
	}

	g, err := game.LoadGameFromRedis(gameId)
	if err != nil {
		http.Error(w, "Error loading game from redis", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	populatePreparedStatsCache(g)
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
