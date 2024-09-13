package app

import (
	"context"
	"sync"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var (
	aiClientInstance *genai.Client
	once             sync.Once
)

type AiModelOption struct {
	Instruction string
	TopK        int32
	TopP        float32
	Temperature float64
}

func GetAIClient(context context.Context) *genai.Client {

	apiKey := GetEnv("AI_API_KEY")

	once.Do(func() {
		aiClientInstance, _ = genai.NewClient(context, option.WithAPIKey(apiKey)) // Inisialisasi client AI
	})
	return aiClientInstance
}

func GetAIModel(client *genai.Client, option *AiModelOption) *genai.GenerativeModel {
	modelID := GetEnv("AI_MODEL_ID")
	model := client.GenerativeModel(modelID)
	model.SetTemperature(*model.Temperature)
	model.SetTopK(option.TopK)
	model.SetTopP(option.TopP)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "text/plain"

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(option.Instruction)},
	}

	return model
}
