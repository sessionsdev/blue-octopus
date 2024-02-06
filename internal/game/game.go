package game

import "log"

type GameMessage interface {
	Provider() string
	Message() string
}

type Game struct {
	GameId             string        `json:"game_id"`
	World              *World        `json:"world"`
	Player             *Player       `json:"player"`
	CentralPlot        string        `json:"central_plot"`
	StoryThreads       []string      `json:"story_threads"`
	SetupMessage       GameMessage   `json:"setup_message"`
	GameMessageHistory []GameMessage `json:"game_message_history"`
	TotalTokensUsed    int           `json:"total_tokens_used"`
}

func (g *Game) BuildCurrentContext(state GameMessage, userPrompt GameMessage) []GameMessage {
	context := []GameMessage{g.SetupMessage, state}

	currentHistory := g.GameMessageHistory
	if len(currentHistory) > 9 {
		// Take the 9 most recent history items
		currentHistory = currentHistory[len(currentHistory)-9:]
	}

	context = append(context, currentHistory...)
	context = append(context, userPrompt)
	// print the current contexts
	for _, message := range context {
		log.Println("-------------------")
		log.Println(message.Message())
	}

	return context
}
