package game

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/sessionsdev/blue-octopus/internal/aiapi"
)

func ProcessGameCommand(command string, gameId string) (string, *Game, error) {
	switch command {
	case "RESET GAME":
		g := InitializeNewGame()
		return fmt.Sprintf("RESET GAME: New game created with id: %s", g.GameId), g, nil
	default:
		g, err := LoadGameFromRedis(gameId)
		if err != nil {
			return "Error loading game. Try the [RESET GAME] command.", nil, err
		}

		narrativeResponse, err := g.processPlayerPrompt(command)
		if err != nil {
			return fmt.Sprintf("An error occured processing the command: %s", command), g, err
		}

		return narrativeResponse, g, nil
	}
}

func (g *Game) processPlayerPrompt(command string) (string, error) {
	messages := []GameMessage{}
	messages = append(messages, GameMessage{Provider: "system", Message: GAME_MASTER_RESPONSABILITY_PROMPT})

	history := g.GetRecentHistory(10)
	messages = append(messages, history...)

	gameState := BuildGameStatePrompt(g)
	gameStateMsg := GameMessage{Provider: "system", Message: gameState}
	messages = append(messages, gameStateMsg)

	userMsg := GameMessage{Provider: "user", Message: command}
	messages = append(messages, userMsg)

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

	g.UpdateGameHistory(userMsg, assistantMessage)
	go g.ReconcileGameState(userMsg, assistantMessage)

	return responseMessage, nil
}

func (g *Game) ReconcileGameState(userMsg GameMessage, assistantMsg GameMessage) {
	messages := []GameMessage{
		{Provider: "system", Message: BuildStateManagerSystemPrompt(g)},
	}

	reconcileStatePrompt := `
		%s

		Reconcile the game state with the following updates:

		Player Action: %s

		Assistant Response: %s
	`

	userMessage := fmt.Sprintf(reconcileStatePrompt, BuildGameStatePrompt(g), userMsg.Message, assistantMsg.Message)
	messages = append(messages, GameMessage{Provider: "user", Message: userMessage})

	// Call the OpenAI Chat API using the client
	response, err := callClient("openai-json", messages)
	if err != nil {
		log.Print("Error calling OpenAI Chat API: ", err)
		return
	}

	// update tokens used
	g.TotalTokensUsed += response.GetTokenUsage()

	// marshal the response message
	responseMessage := response.GetChatCompletion()
	var gameResponse GameStateUpdateResponse
	err = json.Unmarshal([]byte(responseMessage), &gameResponse)
	if err != nil {
		log.Print("Error unmarshaling response: ", err)
		return
	}

	g.UpdateGameState(gameResponse)
	g.SaveGameToRedis()
	g.populatePreparedStatsCache()
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
	case "openai-json":
		aiMessages := []aiapi.AiMessage{}
		for _, message := range messages {
			aiMessages = append(aiMessages, aiapi.AiMessage{
				Provider: message.Provider,
				Message:  message.Message,
			})
		}

		return aiapi.OpenAiJsonClient.DoRequest(aiMessages)
	default:
		return nil, fmt.Errorf("client not found: %s", clientName)
	}
}
