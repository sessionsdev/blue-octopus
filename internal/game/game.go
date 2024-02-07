package game

type GameMessage interface {
	Provider() string
	Message() string
}

type Game struct {
	GameId             string        `json:"game_id"`
	World              *World        `json:"world"`
	Player             *Player       `json:"player"`
	StoryThreads       []string      `json:"story_threads"`
	GameMessageHistory []GameMessage `json:"game_message_history"`
	TotalTokensUsed    int           `json:"total_tokens_used"`
}

func (g *Game) GetRecentHistory(numItems int) []GameMessage {
	currentHistory := g.GameMessageHistory
	if len(currentHistory) > 5 {
		// Take the 5 most recent history items
		return currentHistory[len(currentHistory)-5:]
	} else {
		return currentHistory
	}
}
