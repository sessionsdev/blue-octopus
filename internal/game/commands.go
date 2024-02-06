package game

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/sessionsdev/blue-octopus/internal/aiapi"
)

func (g *Game) ProcessGameCommand(command string) (GameUpdate, error) {
	userMessage := aiapi.OpenAiMessage{Role: "user", Content: command}

	messages := g.BuildCurrentContext(
		getCurrentGameStateMessage(g),
		userMessage,
	)

	// Call the OpenAI Chat API using the client
	response, err := callClient("openai", messages)
	if err != nil {
		return GameUpdate{}, fmt.Errorf("error calling OpenAI Chat API: %w", err)
	}

	// get the raw response message
	responseMessage := response.GetFirstChoiceContent()

	var gameResponse GameUpdate
	err = json.Unmarshal([]byte(responseMessage), &gameResponse)
	if err != nil {
		return GameUpdate{}, fmt.Errorf("error unmarshaling response: %w", err)
	}

	log.Printf("Game Response: %+v", gameResponse)

	assistantMessage := aiapi.OpenAiMessage{Role: "assistant", Content: gameResponse.Response}

	g.UpdateGameState(userMessage, gameResponse, response.Usage.TotalTokens, assistantMessage)

	return gameResponse, nil
}

func callClient(clientName string, messages []GameMessage) (aiapi.OpenAIResponse, error) {
	switch clientName {
	case "openai":
		openAiMessages := make([]aiapi.OpenAiMessage, len(messages))
		for i, message := range messages {
			if message != nil {
				openAiMessages[i] = message.(aiapi.OpenAiMessage)
			}
		}

		return aiapi.Client.CallOpenAIChat(openAiMessages)
	default:
		return aiapi.OpenAIResponse{}, fmt.Errorf("unknown client: %s", clientName)
	}
}

func getCurrentGameStateMessage(game *Game) GameMessage {
	gameState := game.BuildGameStatePromptDetails()

	formattedState := fmt.Sprint(`
	[GAME STATE]
	Current Location: `, gameState.CurrentLocation, `
	Adjacent Locations: `, gameState.AdjacentLocationNames, `
	Inventory: `, gameState.Inventory, `
	Enemies in Location: `, gameState.EnemiesInLocation, `
	Interactive Items: `, gameState.InteractiveItems, `
	Removable Items: `, gameState.RemovableItems, `
	Central Plot: `, gameState.CentralPlot, `
	Story Threads: `, gameState.StoryThreads, `
	`)

	log.Printf("CURRENT GAME STATE: %s", formattedState)

	return aiapi.OpenAiMessage{
		Role:    "system",
		Content: formattedState,
	}
}
