package aiapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var OpenAiClient *OpenAIClient
var OpenAiJsonClient *OpenAIClient

var ModelMap = map[string]string{
	"gpt3": "gpt-3.5-turbo-0125",
	"gpt4": "gpt-4-0125-preview",
}

func init() {
	OpenAiClient = New(ModelMap["gpt3"], 0.7, ResponseFormat{Type: "text"})
	OpenAiJsonClient = New(ModelMap["gpt4"], 0.7, ResponseFormat{Type: "json_object"})
}

func New(model string, temp float64, responseFormat ResponseFormat) *OpenAIClient {
	return &OpenAIClient{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		ClientName:     "openai",
		APIKey:         os.Getenv("OPENAI_API_KEY"),
		Model:          model,
		Temperature:    temp,
		ResponseFormat: responseFormat,
	}
}

// OpenAiMessage represents a single message in the conversation.
type OpenAiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (m *OpenAiMessage) NewMessage(provider string, message string) *OpenAiMessage {
	return &OpenAiMessage{Role: provider, Content: message}
}

// ChatRequest is the request payload for the Chat API.
type ChatRequest struct {
	Model          string          `json:"model"`
	Messages       []OpenAiMessage `json:"messages"`
	Temperature    float64         `json:"temperature"`
	ResponseFormat ResponseFormat  `json:"response_format"`
}

type OpenAiChatResponse struct {
	Completion string `json:"completion"`
	TokensUsed int    `json:"tokens_used"`
}

func (resp *OpenAiChatResponse) GetChatCompletion() string {
	return resp.Completion
}

func (resp *OpenAiChatResponse) GetTokenUsage() int {
	return resp.TokensUsed
}

type ResponseFormat struct {
	Type string `json:"type"`
}

// Choice represents a single choice in the OpenAI response.
type Choice struct {
	Index        int              `json:"index"`
	Message      OpenAiMessage    `json:"message"`
	Logprobs     *json.RawMessage `json:"logprobs"` // Use *json.RawMessage for nullability
	FinishReason string           `json:"finish_reason"`
}

// Usage details the token usage of the OpenAI response.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// OpenAIResponse represents the structure of the response from OpenAI Chat API.
type OpenAIResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
}

func (resp *OpenAIResponse) GetChatCompletion() string {
	return resp.Choices[0].Message.Content
}

func (resp *OpenAIResponse) GetTokenUsage() int {
	return resp.Usage.TotalTokens
}

type OpenAIClient struct {
	Client         *http.Client
	ClientName     string
	APIKey         string
	Model          string
	Temperature    float64
	ResponseFormat ResponseFormat
}

func (c *OpenAIClient) DoRequest(userMessages []AiMessage) (*OpenAiChatResponse, error) {
	openAiMessage := convertMessageType(userMessages)

	chatRequest := ChatRequest{
		Model:          c.Model,
		Messages:       openAiMessage,
		Temperature:    c.Temperature,
		ResponseFormat: c.ResponseFormat,
	}

	requestBody, err := json.Marshal(chatRequest)
	if err != nil {
		log.Panicf("HERE I AM!")
		return &OpenAiChatResponse{}, fmt.Errorf("error marshaling request: %w", err)
	}

	url := "https://api.openai.com/v1/chat/completions"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return &OpenAiChatResponse{}, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	// Log request attempt
	resp, err := c.Client.Do(req)
	if err != nil {
		return &OpenAiChatResponse{}, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// print response body
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		log.Printf("response body: %s", buf.String())
		return &OpenAiChatResponse{}, fmt.Errorf("received non-OK response status: %s", resp.Status)
	}

	var openAiResponse OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAiResponse); err != nil {
		return &OpenAiChatResponse{}, fmt.Errorf("error decoding response: %w", err)
	}

	// prettyPrint(openAiResponse)
	log.Println("CLIENT RESPONSE TYPE: ", c.ResponseFormat.Type)
	log.Println("OpenAI Response: ", openAiResponse.Choices[0].Message.Content)

	return &OpenAiChatResponse{
		Completion: openAiResponse.GetChatCompletion(),
		TokensUsed: openAiResponse.GetTokenUsage(),
	}, nil
}

func convertMessageType(messages []AiMessage) []OpenAiMessage {
	var openAiMessages []OpenAiMessage
	for _, m := range messages {
		openAiMessages = append(openAiMessages, OpenAiMessage{Role: m.Provider, Content: m.Message})
	}
	return openAiMessages
}
