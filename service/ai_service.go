package service

import (
	"adiubaidah/adi-bot/model/ai"
	"adiubaidah/adi-bot/model/history"
	"context"

	"github.com/google/generative-ai-go/genai"
)

type AiService interface {
	CreateModel(ctx context.Context, request ai.AiRequest) *genai.GenerativeModel
	GenerateResponse(ctx context.Context, model *genai.GenerativeModel, histories *[]history.History, input string) (string, error)
}
