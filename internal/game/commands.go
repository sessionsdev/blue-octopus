package game

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/sessionsdev/blue-octopus/internal/aiapi"
)

func (g *Game) ProcessGameCommand(command string) (string, error) {
	userMessage := aiapi.OpenAiMessage{Role: "user", Content: command}
	gameMasterMessage := aiapi.OpenAiMessage{Role: "system", Content: GAME_MASTER_RESPONSABILITY_PROMPT}
	responseProtocal := aiapi.OpenAiMessage{Role: "system", Content: GAME_MASTER_RESPONSE_PROTOCOL_PROMPT}
	gameState := g.BuildGameStateDetails()
	formattedState := fmt.Sprintf(
		GAME_STATE_PROMPT,
		gameState.CurrentLocation,
		gameState.AdjacentLocationNames,
		gameState.Inventory,
		gameState.InteractiveItems,
		gameState.Obstacles,
		gameState.StoryThreads,
	)

	gameMasterStateMessage := aiapi.OpenAiMessage{Role: "system", Content: formattedState}

	history := g.GetRecentHistory(5)
	messages := []GameMessage{}
	messages = append(messages, gameMasterMessage, responseProtocal, gameMasterStateMessage)
	messages = append(messages, history...)
	messages = append(messages, userMessage)

	// Call the OpenAI Chat API using the client
	response, err := callClient("openai", messages)
	if err != nil {
		return "", fmt.Errorf("error calling OpenAI Chat API: %w", err)
	}

	// update tokens used
	g.TotalTokensUsed += response.Usage.TotalTokens

	// get the raw response message
	responseMessage := response.GetFirstChoiceContent()
	log.Print("Response Message: ", responseMessage)

	assistantMessage := aiapi.OpenAiMessage{Role: "assistant", Content: responseMessage}

	g.UpdateGameHistory(userMessage, assistantMessage)
	g.reconcileGameState(responseMessage)

	return responseMessage, nil
}

func (g *Game) reconcileGameState(narrativeUpdate string) {
	stateManagerResponsibilityMessage := aiapi.OpenAiMessage{Role: "system", Content: STATE_MANAGER_RESPONSABILITY_PROMPT}
	stateManagerResponseProtocol := aiapi.OpenAiMessage{Role: "system", Content: STATE_MANAGER_RESPONSE_PROTOCOL_PROMPT}
	gameState := g.BuildGameStateDetails()
	formattedState := fmt.Sprintf(
		GAME_STATE_PROMPT,
		gameState.CurrentLocation,
		gameState.AdjacentLocationNames,
		gameState.Inventory,
		gameState.InteractiveItems,
		gameState.Obstacles,
		gameState.StoryThreads,
	)

	userPromptString := fmt.Sprint(`
	[NARRATIVE UPDATE]
	`, narrativeUpdate, `
	`)

	gameMasterStateMessage := aiapi.OpenAiMessage{Role: "system", Content: formattedState}

	messages := []GameMessage{
		stateManagerResponsibilityMessage,
		stateManagerResponseProtocol,
		gameMasterStateMessage,
	}

	history := g.GetRecentHistory(5)
	messages = append(messages, history...)
	messages = append(messages, aiapi.OpenAiMessage{Role: "user", Content: userPromptString})

	// Call the OpenAI Chat API using the client
	response, err := callClient("openai", messages)
	if err != nil {
		log.Print("Error calling OpenAI Chat API: ", err)
		return
	}

	// update tokens used
	g.TotalTokensUsed += response.Usage.TotalTokens

	// marshal the response message
	responseMessage := response.GetFirstChoiceContent()
	log.Print("Response Message: ", responseMessage)

	var gameResponse GameStateDetails
	err = json.Unmarshal([]byte(responseMessage), &gameResponse)
	if err != nil {
		log.Print("Error unmarshaling response: ", err)
		return
	}

	g.UpdateGameState(gameResponse)
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
