package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sessionsdev/blue-octopus/internal/openai"
)

type Command struct {
	Command string `json:"command"`
}

type GameStateSystemPrompt struct {
	Context string    `json:"context"`
	State   GameState `json:"gameState"`
}

type GameState struct {
	CharacterName     string   `json:"characterName"`
	CurrentLocation   string   `json:"currentLocation"`
	Inventory         []string `json:"inventory"`
	PreviousLocations []string `json:"previousLocations"`
	EnemiesInLocation []string `json:"enemiesInLocation"`
}

type OpenAiGameResponse struct {
	Response             string    `json:"response"`
	ProposedStateChanges GameState `json:"proposedStateChanges"`
}

func (gs *GameState) JsonString() string {
	jsonData, err := json.Marshal(gs)
	if err != nil {
		return ""
	}
	return "gameState: " + string(jsonData)
}

var state = GameState{
	CharacterName:     "Bob",
	CurrentLocation:   "The Forest",
	Inventory:         []string{"sword", "shield"},
	PreviousLocations: []string{"The Town", "The Tavern"},
	EnemiesInLocation: []string{"goblin", "troll"},
}

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

func processGameCommand(command string) (OpenAiGameResponse, error) {
	gameMessages := []openai.Message{
		{
			Role: "system",
			Content: `You are the narrator of a text-based adventure game. ` +
				`You are in control of the game's story and the world in which it takes place. ` +
				`You can create new locations, enemies, and items, and you can describe the game's world to the player ` +
				`and present puzzles, riddles and challenges, You can also respond to the player's commands and questions. ` +
				`The game's state will be passed to each time along with the players command.`,
		},
		{
			Role: "system",
			Content: `Always respond in the format a json object 'response:' key with the chat completion, ` +
				`and a 'proposeStateChanges:' key with nested json object of the proposed state changes with the following schema: ` +
				`
				{
					"characterName": "Character's Name",
					"currentLocation": "Current Location",
					"inventory": ["item1", "item2", "item3"],
					"previousLocations": ["location1", "location2"],
					"enemiesInLocation": ["enemy1", "enemy2"]
				}
				`,
		},
		{
			Role:    "system",
			Content: state.JsonString(),
		},
	}

	gameMessages = append(gameMessages, gameMessageHistory...)

	userMessage := openai.Message{Role: "user", Content: command}
	gameMessages = append(gameMessages, userMessage)

	// print the game messages with role
	for _, message := range gameMessages {
		fmt.Printf("%s: %s\n\n", message.Role, message.Content)
	}

	// Call the OpenAI Chat API using the client
	response, err := client.CallOpenAIChat(gameMessages)
	if err != nil {
		return OpenAiGameResponse{}, fmt.Errorf("error calling OpenAI Chat API: %w", err)
	}

	// print full response
	fmt.Printf("Full Response: %+v\n", response)

	// get the raw response message
	responseMessage := response.GetFirstChoice()

	// turn responseMessage into a OpenAiGameResponse
	var gameResponse OpenAiGameResponse
	err = json.Unmarshal([]byte(responseMessage), &gameResponse)
	if err != nil {
		return OpenAiGameResponse{}, fmt.Errorf("error unmarshaling response: %w", err)
	}

	assistentResponse := openai.Message{Role: "assistant", Content: gameResponse.Response}
	proposedStateChanges := gameResponse.ProposedStateChanges
	//print the proposed state changes
	fmt.Printf("Proposed State Changes: %+v\n", proposedStateChanges)

	gameMessageHistory = append(gameMessageHistory, userMessage)
	gameMessageHistory = append(gameMessageHistory, assistentResponse)

	return gameResponse, nil
}

func resolveStateChanges(proposedStateChanges GameState) {
	// update the game state
	state.CurrentLocation = proposedStateChanges.CurrentLocation
	state.Inventory = proposedStateChanges.Inventory
	state.PreviousLocations = proposedStateChanges.PreviousLocations
	state.EnemiesInLocation = proposedStateChanges.EnemiesInLocation
}
