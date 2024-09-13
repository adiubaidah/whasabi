package service

import (
	"adiubaidah/adi-bot/app"
	"adiubaidah/adi-bot/helper"
	"adiubaidah/adi-bot/model/ai"
	"adiubaidah/adi-bot/model/history"
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/generative-ai-go/genai"
)

type AiServiceImpl struct {
	Client   *genai.Client
	Validate *validator.Validate
}

func NewAiService(client *genai.Client, validate *validator.Validate) AiService {
	return &AiServiceImpl{
		Client:   client,
		Validate: validate,
	}
}

func (service *AiServiceImpl) CreateModel(ctx context.Context, request ai.AiRequest) *genai.GenerativeModel {

	err := service.Validate.Struct(request)
	helper.PanicIfError("", err)

	// Get the model
	option := app.AiModelOption{
		Instruction: request.Instruction,
		TopK:        request.TopK,
		TopP:        request.TopP,
		Temperature: request.Temperature,
	}

	model := app.GetAIModel(service.Client, &option)
	return model
}

func (service *AiServiceImpl) GenerateResponse(ctx context.Context, model *genai.GenerativeModel, histories *[]history.History, input string) (string, error) {
	var sessionHistory []*genai.Content
	for _, history := range *histories {
		sessionHistory = append(sessionHistory, &genai.Content{
			Parts: []genai.Part{genai.Text(history.Content)},
		})
	}

	// Generate response
	session := model.StartChat()
	session.History = sessionHistory

	resp, err := session.SendMessage(ctx, genai.Text(input))
	helper.PanicIfError("Error saat mengambil respon:", err)

	return string(resp.Candidates[0].Content.Parts[0].(genai.Text)), nil
}
