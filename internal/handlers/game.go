package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sessionsdev/blue-octopus/internal/game"
	"github.com/sessionsdev/blue-octopus/internal/openai"
)

type Command struct {
	Command string `json:"command"`
}

type UserPromptWithState struct {
	Prompt               string `json:"prompt"`
	ProposedStateChanges string `json:"proposed_state_changes"`
}

var Game *game.Game

func init() {
	// Initialize a new game
	Game = game.InitializeNewGame(openai.Message{Role: "system", Content: game.SETUP_PROMPT})
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

	var html string

	switch command {
	case "RESET GAME":
		Game = game.InitializeNewGame(openai.Message{Role: "system", Content: game.SETUP_PROMPT})
		html = "<p>============ GAME RESET ============</p>"
	default:
		gameUpdate, err := Game.ProcessGameCommand(command)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		html = fmt.Sprintf("<p>PLAYER: %s\n<p>NARRATOR: %s</p>", command, gameUpdate.Response)
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func HandleGameState(w http.ResponseWriter, r *http.Request) {
	// For GET requests, return the current full game state
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		newGameJson, err := json.Marshal(Game.GetGameJson())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(newGameJson)
		return
	}
}

func HandleWorldTreeVisualization(w http.ResponseWriter, r *http.Request) {
	// For GET requests, return the current full game state
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/plain")
		Game.World.VisualizeLocationTree()
		return
	}
}
