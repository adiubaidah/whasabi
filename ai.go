package main

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func GetAIModel() (*genai.GenerativeModel, *genai.Client) {
	ctx := context.Background()

	// Retrieve API key and model ID
	apiKey, err := GetConfig("GENAI_API_KEY")
	PanicIfError("Failed to get API key", err)

	modelID, err := GetConfig("GENAI_MODEL_ID")
	PanicIfError("Failed to get Model ID", err)

	systemInstruction, err := GetConfig("GENAI_SYSTEM_INSTRUCTION")
	PanicIfError("Error retrieving system instruction:", err)

	// Initialize the AI client
	aiClient, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	PanicIfError("Failed to initialize AI client", err)

	// Set model parameters
	model := aiClient.GenerativeModel(modelID)
	model.SetTemperature(0.9)
	model.SetTopK(64)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "text/plain"

	// Set system instructions
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemInstruction)},
	}

	// Return both the model and the AI client (so the caller can close it when done)
	return model, aiClient
}
