package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sessionsdev/blue-octopus/internal/openai"
)

func ServeHelloWorldAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<div> Hello World!</div>")
}

func GenerateNewText(w http.ResponseWriter, r *http.Request) {
	// Replace these with your actual API key and the OpenAI API URL
	apiKey := "XXXXX"
	apiUrl := "https://api.openai.com/v1/chat/completions"

	// Create a new OpenAI client
	client := openai.New(apiKey, apiUrl)

	// Define your model and user message
	model := "gpt-3.5-turbo"
	userMessage := "Why is the sky blue?"
	temperature := 0.7 // Adjust the temperature as needed

	// Call the OpenAI Chat API using the client
	response, err := client.CallOpenAIChat(model, userMessage, temperature)
	if err != nil {
		log.Fatalf("Error calling OpenAI Chat: %v", err)
	}

	// Print the full response
	fmt.Printf("Message: %+v\n", response.Choices[0].Message.Content)
}
