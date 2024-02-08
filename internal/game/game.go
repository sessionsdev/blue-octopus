package game

type Message interface {
	NewMessage(string, string) Message
}

type GameMessage struct {
	Provider string `json:"provider"`
	Message  string `json:"message"`
}

func (m *GameMessage) NewMessage(provider string, message string) Message {
	return &GameMessage{Provider: provider, Message: message}
}

type Game struct {
	GameId             string        `json:"game_id"`
	World              *World        `json:"world"`
	Player             *Player       `json:"player"`
	MainQuest          string        `json:"main_quest"`
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
