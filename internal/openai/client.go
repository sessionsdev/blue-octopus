package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Message represents a single message in the conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest is the request payload for the Chat API.
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

// Choice represents a single choice in the OpenAI response.
type Choice struct {
	Index        int              `json:"index"`
	Message      Message          `json:"message"`
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

type OpenAIClient struct {
	Client *http.Client
	APIKey string
	URL    string
}

func New(apiKey, url string) *OpenAIClient {
	return &OpenAIClient{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		APIKey: apiKey,
		URL:    url,
	}
}

func (c *OpenAIClient) CallOpenAIChat(model, userMessage string, temperature float64) (OpenAIResponse, error) {
	chatRequest := ChatRequest{
		Model:       model,
		Messages:    []Message{{Role: "user", Content: userMessage}},
		Temperature: temperature,
	}

	requestBody, err := json.Marshal(chatRequest)
	if err != nil {
		return OpenAIResponse{}, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(requestBody))
	if err != nil {
		return OpenAIResponse{}, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.Client.Do(req)
	if err != nil {
		return OpenAIResponse{}, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return OpenAIResponse{}, fmt.Errorf("received non-OK response status: %s", resp.Status)
	}

	var result OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return OpenAIResponse{}, fmt.Errorf("error decoding response: %w", err)
	}

	return result, nil
}
