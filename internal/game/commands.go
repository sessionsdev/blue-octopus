package game

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/sessionsdev/blue-octopus/internal/openai"
)

func (newGame *Game) ProcessGameCommand(command string) (GameUpdate, error) {
	userMessage := openai.Message{Role: "user", Content: command}

	messages := newGame.BuildCurrentContext(
		getCurrentGameStateMessage(newGame),
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

	assistantMessage := openai.Message{Role: "assistant", Content: gameResponse.Response}

	newGame.UpdateGameState(userMessage, gameResponse, response.Usage.TotalTokens, assistantMessage)

	return gameResponse, nil
}

func callClient(clientName string, messages []GameMessage) (openai.OpenAIResponse, error) {
	switch clientName {
	case "openai":
		openAiMessages := make([]openai.Message, len(messages))
		for i, message := range messages {
			if message != nil {
				openAiMessages[i] = message.(openai.Message)
			}
		}

		return openai.Client.CallOpenAIChat(openAiMessages)
	default:
		return openai.OpenAIResponse{}, fmt.Errorf("unknown client: %s", clientName)
	}
}

func getCurrentGameStateMessage(game *Game) GameMessage {
	content := game.BuildGameStatePromptDetails()

	return openai.Message{
		Role:    "system",
		Content: content.GetJsonOrEmptyString(),
	}
}
