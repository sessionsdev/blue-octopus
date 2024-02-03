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

type GameStatePromptDetails struct {
	CurrentLocation   *game.Location `json:"current_location"`
	Inventory         []string       `json:"inventory"`
	EnemiesInLocation []string       `json:"enemies_in_location"`
}

type OpenAiGameResponse struct {
	Response             string                    `json:"response"`
	ProposedStateChanges game.ProposedStateChanges `json:"proposed_state_changes"`
}

var newGame = game.InitializeNewGame()

var gameMessageHistory = []openai.Message{}

var client = openai.New("gpt-4-0125-preview", 0.7, openai.ResponseFormat{Type: "json_object"})

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

	gameUpdate, err := processGameCommand(command)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<p>" + gameUpdate.Response + "</p>"))
}

func HandleGameState(w http.ResponseWriter, r *http.Request) {
	// For GET requests, return the current full game state
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		newGameJson, err := json.Marshal(newGame)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(newGameJson)
		return
	}
}

func getPrefixMessages() []openai.Message {
	return []openai.Message{
		{
			Role: "system",
			Content: `You are the game master of a text-based adventure game. ` +
				`You are in control of the game's story and the world in which it takes place. ` +
				`You can create new locations, enemies, and items, and you can describe the game's world to the player ` +
				`and present puzzles, riddles and challenges, You can also respond to the player's commands and questions. `,
		},
		{
			Role: "system",
			Content: `Always respond in the format a json object in the following format:, ` +
				`
				{
					"response": "You enter the Small Room. You see a chest and a door, and a sword.  There is a shadowy figure What would you like to do?"
					"proposed_state_changes": {
						"current_location_name": "Small Room",
						"interactive_objects_in_location": ["chest", "door"],
						"removable_items": ["sword"],
						"updated_enemies_in_location": ["shadowy figure"],
						"player_inventory": ["potion"]
					}
				}

				Another example:
				{
					"response": "Opened the chest and take the sword inside before dispatching the shadowy figure. The door opens. What would you like to do next?"
					"proposed_state_changes": {
						"current_location_name": "Small Room",
						"interactive_objects_in_location": ["chest", "door"],
						"updated_removable_items": [],
						"enemies_in_location": [],
						"player_inventory": ["potion", "sword"]
					}
				}

				Another example:
				{
					"response": "You go through the door.  On the other side of the door is the Library.  There's a book on the table. What would you like to do next?"
					"proposed_state_changes": {
						"current_location_name": "Library",
						"interactive_objects_in_location": ["book"],
						"removable_items": [],
						"enemies_in_location": [],
						"player_inventory": ["potion", "sword"]
					}
				}
				`,
		},
	}
}

func getCurrentGameStateMessage() openai.Message {
	return openai.Message{
		Role:    "system",
		Content: newGame.GetPartialJsonRepresentation(),
	}
}

func processGameCommand(command string) (OpenAiGameResponse, error) {
	gameMessages := getPrefixMessages()
	gameMessages = append(gameMessages, getCurrentGameStateMessage())
	gameMessages = append(gameMessages, gameMessageHistory...)

	userMessage := openai.Message{Role: "user", Content: command}
	gameMessages = append(gameMessages, userMessage)

	// Call the OpenAI Chat API using the client
	response, err := client.CallOpenAIChat(gameMessages)
	if err != nil {
		return OpenAiGameResponse{}, fmt.Errorf("error calling OpenAI Chat API: %w", err)
	}

	// get the raw response message
	responseMessage := response.GetFirstChoice()

	// print the response message
	fmt.Printf("Response: %s\n\n", responseMessage)

	// turn responseMessage into a OpenAiGameResponse
	var gameResponse OpenAiGameResponse
	err = json.Unmarshal([]byte(responseMessage), &gameResponse)
	if err != nil {
		return OpenAiGameResponse{}, fmt.Errorf("error unmarshaling response: %w", err)
	}

	assistentResponse := openai.Message{Role: "assistant", Content: gameResponse.Response}
	proposedStateChanges := gameResponse.ProposedStateChanges

	// print the proposed state changes
	fmt.Printf("Proposed State Changes: %+v\n\n", proposedStateChanges)

	gameMessageHistory = append(gameMessageHistory, userMessage)
	gameMessageHistory = append(gameMessageHistory, assistentResponse)
	newGame.UpdateGameState(proposedStateChanges, response.Usage.TotalTokens)

	return gameResponse, nil
}
