package game

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/sessionsdev/blue-octopus/internal/aiapi"
)

func (g *Game) ProcessGameCommand(command string) (string, error) {
	switch command {
	case "RESET GAME":
		g := InitializeNewGame()
		return fmt.Sprintf("RESET GAME: New game created with id: %s", g.GameId), nil
	default:
		narrativeResponse, err := g.processPlayerPrompt(command)
		if err != nil {
			return fmt.Sprintf("An error occured processing the command: %s", command), err
		}

		return narrativeResponse, nil
	}
}

func (g *Game) processPlayerPrompt(command string) (string, error) {
	promptWithState := BuildGameMasterStatePrompt(g, command)
	userMessage := GameMessage{Provider: "user", Message: promptWithState}

	history := g.GetRecentHistory(5)
	messages := []GameMessage{}
	messages = append(messages, GameMessage{Provider: "system", Message: GAME_MASTER_RESPONSABILITY_PROMPT})
	messages = append(messages, history...)
	messages = append(messages, userMessage)

	// Call the OpenAI Chat API using the client
	response, err := callClient("openai", messages)
	if err != nil {
		return "", fmt.Errorf("error calling OpenAI Chat API: %w", err)
	}

	// update tokens used
	g.TotalTokensUsed += response.GetTokenUsage()

	// get the raw response message
	responseMessage := response.GetChatCompletion()

	assistantMessage := GameMessage{Provider: "assistant", Message: responseMessage}

	g.UpdateGameHistory(userMessage, assistantMessage)
	go g.reconcileGameState(responseMessage)

	return responseMessage, nil
}

func (g *Game) reconcileGameState(narrativeUpdate string) {
	stateManagerResponsibilityMessage := GameMessage{Provider: "system", Message: STATE_MANAGER_RESPONSABILITY_PROMPT}
	stateManagerResponseProtocol := GameMessage{Provider: "system", Message: STATE_MANAGER_RESPONSE_PROTOCOL_PROMPT}
	gameState := g.BuildGameStateDetails()
	formattedState := fmt.Sprintf(
		GAME_STATE_PROMPT,
		gameState.CurrentLocation,
		gameState.Inventory,
		gameState.InteractiveItems,
		gameState.Obstacles,
		gameState.StoryThreads,
	)

	gameMasterStateMessage := GameMessage{Provider: "system", Message: formattedState}

	messages := []GameMessage{
		stateManagerResponsibilityMessage,
		stateManagerResponseProtocol,
		gameMasterStateMessage,
	}

	history := g.GetRecentHistory(3)
	messages = append(messages, history...)

	userPromptString := fmt.Sprint(`
	[NARRATIVE UPDATE]
	`, narrativeUpdate, `
	`)
	userPromptMessage := GameMessage{Provider: "user", Message: userPromptString}
	messages = append(messages, userPromptMessage)

	// Call the OpenAI Chat API using the client
	response, err := callClient("openai", messages)
	if err != nil {
		log.Print("Error calling OpenAI Chat API: ", err)
		return
	}

	// update tokens used
	g.TotalTokensUsed += response.GetTokenUsage()

	// marshal the response message
	responseMessage := response.GetChatCompletion()
	var gameResponse GameStateDetails
	err = json.Unmarshal([]byte(responseMessage), &gameResponse)
	if err != nil {
		log.Print("Error unmarshaling response: ", err)
		return
	}

	g.UpdateGameState(gameResponse)
}

func callClient(clientName string, messages []GameMessage) (aiapi.ChatResponse, error) {
	switch clientName {
	case "openai":
		aiMessages := []aiapi.AiMessage{}
		for _, message := range messages {
			aiMessages = append(aiMessages, aiapi.AiMessage{
				Provider: message.Provider,
				Message:  message.Message,
			})
		}

		return aiapi.OpenAiClient.DoRequest(aiMessages)
	default:
		return nil, fmt.Errorf("client not found: %s", clientName)
	}
}
