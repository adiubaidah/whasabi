package service

import (
	"adiubaidah/adi-bot/model"
	"context"

	"github.com/google/generative-ai-go/genai"
)

type AiService interface {
	CreateConfiguration(ctx context.Context, configuration model.AiConfiguration) *model.Ai
	GetConfiguration(ctx context.Context) *model.Ai
	GetModel(ctx context.Context, request model.AiConfiguration) *genai.GenerativeModel
	GenerateResponse(ctx context.Context, model *genai.GenerativeModel, histories *[]model.History, input string) (string, error)
}
