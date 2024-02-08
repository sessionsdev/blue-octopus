package aiapi

type ChatResponse interface {
	GetChatCompletion() string
	GetTokenUsage() int
}

type AiMessage struct {
	Provider string `json:"provider"`
	Message  string `json:"message"`
}

func (m *AiMessage) NewMessage(provider string, message string) *AiMessage {
	return &AiMessage{Provider: provider, Message: message}
}

type AIClient interface {
	DoRequest(messages []AiMessage) (*ChatResponse, error)
}
