package game

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/sessionsdev/blue-octopus/internal/aiapi"
)

var gameCommandProcessing bool

// func ProcessGameCommand(ctx context.Context, command string) (string, error) {
// 	userContext := ctx.Value("user")
// 	if userContext == nil {
// 		return "Unauthorized", fmt.Errorf("Unauthorized")
// 	}

// 	user := userContext.(*auth.User)
// 	username := user.Email

// 	switch command {
// 	case "RESET GAME":
// 		g := InitializeNewGame()
// 		SaveGameToRedis(ctx, g, username)
// 		return fmt.Sprintf("RESET GAME: New game created!"), nil
// 	default:
// 		g, err := LoadGameFromRedis(ctx, username)
// 		if err != nil {
// 			log.Println("Error loading game from redis: ", err)
// 			return `No game found. Try using the "RESET GAME" command`, nil
// 		}

// 		narrativeResponse, err := g.ProcessPlayerPrompt(ctx, command)
// 		if err != nil {
// 			return fmt.Sprintf("An error occured processing the command: %s", command), err
// 		}

// 		return narrativeResponse, nil
// 	}
// }

func (g *Game) ProcessPlayerPrompt(command string) (string, error) {
	messages := []GameMessage{
		{Provider: "system", Message: GAME_MASTER_RESPONSABILITY_PROMPT},
		{Provider: "system", Message: BuildGameMasterStatePrompt(g)},
	}

	history := g.GetRecentHistory(20)
	messages = append(messages, history...)
	messages = append(messages, GameMessage{Provider: "user", Message: command})

	// Call the OpenAI Chat API using the client
	response, err := callClient("openai", messages)
	if err != nil {
		return "", fmt.Errorf("error calling OpenAI Chat API: %w", err)
	}

	// update tokens used
	g.TotalTokensUsed += response.GetTokenUsage()

	// get the raw response message
	responseMessage := response.GetChatCompletion()

	userMessage := GameMessage{Provider: "user", Message: command}
	assistantMessage := GameMessage{Provider: "assistant", Message: responseMessage}

	g.UpdateGameHistory(userMessage, assistantMessage)

	return responseMessage, nil
}

func (g *Game) reconcileGameState() {
	messages := []GameMessage{
		{Provider: "system", Message: STATE_MANAGER_RESPONSE_PROTOCOL_PROMPT},
		{Provider: "system", Message: BuildStateManagerPrompt(g)},
	}

	messages = append(messages, g.GetRecentHistory(5)...)

	reconcileStatePrompt := `Reconcile the game state with the previous messages and respond with a structured JSON object.`
	messages = append(messages, GameMessage{Provider: "user", Message: reconcileStatePrompt})

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
	var gameStateResponse GameStateUpdateResponse
	err = json.Unmarshal([]byte(responseMessage), &gameStateResponse)
	if err != nil {
		log.Print("Error unmarshaling response: ", err)
		return
	}

	g.UpdateGameState(gameStateResponse)
}

type StoryThreadsResponse struct {
	StoryThreads []string `json:"story_threads"`
}

func (g *Game) progressStoryThreads() {
	mostRecentAssistantMessage := g.GetRecentHistory(1)[0]
	userMsg := g.GetRecentHistory(2)[1]

	var userMessage string
	if len(g.StoryThreads) > 10 {
		userMessage = BuildProgressiveSummaryPrompt(g.StoryThreads)
	} else {
		userMessage = BuildGameSummaryCurrentStatePrompt(g.StoryThreads, userMsg.Message, mostRecentAssistantMessage.Message)
	}

	messages := []GameMessage{
		{Provider: "system", Message: GAME_SUMMARY_MANAGER_PROMPT},
		{Provider: "user", Message: userMessage},
	}

	response, err := callClient("openai-json", messages)
	if err != nil {
		log.Print("Error calling OpenAI Chat API: ", err)
		return
	}

	g.TotalTokensUsed += response.GetTokenUsage()

	responseMessage := response.GetChatCompletion()
	var storyThreadsResponse StoryThreadsResponse
	err = json.Unmarshal([]byte(responseMessage), &storyThreadsResponse)
	if err != nil {
		log.Print("Error unmarshaling response: ", err)
		return
	}

	g.StoryThreads = storyThreadsResponse.StoryThreads
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
