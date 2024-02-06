package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/sessionsdev/blue-octopus/internal/aiapi"
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
	Enemies   []string
}

var PreparedStatsCache *PreparedStats = &PreparedStats{}

func populatePreparedStatsCache(g *game.Game) {
	PreparedStatsCache.Location = g.World.CurrentLocation.LocationName
	PreparedStatsCache.Inventory = g.Player.Inventory
	PreparedStatsCache.Enemies = g.World.CurrentLocation.EnemiesInLocation

	if len(PreparedStatsCache.Inventory) == 0 {
		PreparedStatsCache.Inventory = []string{"You're not carrying anything."}
	}

	if len(PreparedStatsCache.Enemies) == 0 {
		PreparedStatsCache.Enemies = []string{"You are safe."}
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
		log.Println("No game id cookie found, creating a new game")

		g = game.InitializeNewGame(aiapi.OpenAiMessage{Role: "system", Content: game.SETUP_PROMPT})
		g.SaveGameToRedis()
		log.Printf("New game created with id: %s", g.GameId)
	} else {
		loadedGame, err := game.LoadGameFromRedis(gameIdFromCookie)
		if err != nil {
			log.Println("Error loading game from redis: ", err)
			g = game.InitializeNewGame(aiapi.OpenAiMessage{Role: "system", Content: game.SETUP_PROMPT})
			g.SaveGameToRedis()
		}
		g = loadedGame
	}

	setGameIdCookie(w, g.GameId)

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

	// get game id cookie
	cookie, err := r.Cookie("GameId")
	if err != nil {
		http.Error(w, "No game id cookie found", http.StatusBadRequest)
		return
	}

	g, err := game.LoadGameFromRedis(cookie.Value)
	if err != nil || g == nil {
		if err != nil {
			w.Write([]byte("Error loading game: " + err.Error()))
		} else {
			w.Write([]byte("No game found for id: " + cookie.Value))
		}
		return
	}

	var html string

	switch command {
	case "RESET GAME":
		g = game.InitializeNewGame(aiapi.OpenAiMessage{Role: "system", Content: game.SETUP_PROMPT})
		setGameIdCookie(w, g.GameId)
		html = fmt.Sprintf("RESET GAME: New game created with id: %s", g.GameId)
	default:
		gameUpdate, err := g.ProcessGameCommand(command)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		html = fmt.Sprintf("<p>PLAYER: %s\n<p>GAME MASTER: %s</p>", command, gameUpdate.Response)
	}

	populatePreparedStatsCache(g)
	g.SaveGameToRedis()
	setGameIdCookie(w, g.GameId)
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Trigger-After-Settle", "stats-update")
	w.Write([]byte(html))
}

func ServeGameStats(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/stats-panel.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, PreparedStatsCache)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	clearPreparedStatsCache()
	w.Header().Set("Content-Type", "text/html")
}

func HandleGameState(w http.ResponseWriter, r *http.Request) {
	// For GET requests, return the current full game state
	if r.Method == http.MethodGet {
		// get game id cookie
		cookie, err := r.Cookie("GameId")
		if err != nil {
			http.Error(w, "No game id cookie found", http.StatusBadRequest)
			return
		}

		g, err := game.LoadGameFromRedis(cookie.Value)
		if err != nil || g == nil {
			if err != nil {
				w.Write([]byte("Error loading game: " + err.Error()))
			} else {
				w.Write([]byte("No game found for id: " + cookie.Value))
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		newGameJson, err := json.Marshal(g)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Write(newGameJson)
	}
}

func HandleWorldTreeVisualization(w http.ResponseWriter, r *http.Request) {
	// For GET requests, return the current full game state
	if r.Method == http.MethodGet {
		// get game id cookie
		cookie, err := r.Cookie("GameId")
		if err != nil {
			http.Error(w, "No game id cookie found", http.StatusBadRequest)
			return
		}

		g, err := game.LoadGameFromRedis(cookie.Value)
		if err != nil {
			w.Write([]byte("No game found for id: " + cookie.Value))
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Location Tree Visualization\n"))
		w.Write([]byte(g.World.VisualizeLocationTree()))
	}
}
