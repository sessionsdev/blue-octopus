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
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
	}

	if game.PreparedStatsCache == nil {
		http.Error(w, "No stats available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	executeTemplate(w, "templates/stats-panel.html", "stats-panel", game.PreparedStatsCache)
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
