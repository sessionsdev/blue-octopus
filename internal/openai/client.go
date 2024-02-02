package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Message represents a single message in the conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest is the request payload for the Chat API.
type ChatRequest struct {
	Model          string         `json:"model"`
	Messages       []Message      `json:"messages"`
	Temperature    float64        `json:"temperature"`
	ResponseFormat ResponseFormat `json:"response_format"`
}

type ResponseFormat struct {
	Type string `json:"type"`
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

func (resp *OpenAIResponse) GetFirstChoice() string {
	return resp.Choices[0].Message.Content
}

type OpenAIClient struct {
	Client         *http.Client
	APIKey         string
	Model          string
	Temperature    float64
	ResponseFormat ResponseFormat
}

func New(model string, temp float64, responseFormat ResponseFormat) *OpenAIClient {
	return &OpenAIClient{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		APIKey:         os.Getenv("OPENAI_API_KEY"),
		Model:          model,
		Temperature:    temp,
		ResponseFormat: responseFormat,
	}
}

func (c *OpenAIClient) CallOpenAIChat(userMessages []Message) (OpenAIResponse, error) {
	chatRequest := ChatRequest{
		Model:          "gpt-4-0125-preview",
		Messages:       userMessages,
		Temperature:    0.7,
		ResponseFormat: ResponseFormat{Type: "json_object"},
	}

	requestBody, err := json.Marshal(chatRequest)
	if err != nil {
		return OpenAIResponse{}, fmt.Errorf("error marshaling request: %w", err)
	}

	url := "https://api.openai.com/v1/chat/completions"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
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
		// print response body
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		log.Printf("response body: %s", buf.String())
		return OpenAIResponse{}, fmt.Errorf("received non-OK response status: %s", resp.Status)
	}

	var openAiResponse OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAiResponse); err != nil {
		return OpenAIResponse{}, fmt.Errorf("error decoding response: %w", err)
	}

	// print total tokens used
	fmt.Printf("Total tokens used: %d\n", openAiResponse.Usage.TotalTokens)

	return openAiResponse, nil
}
