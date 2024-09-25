package service

import (
	"adiubaidah/adi-bot/model"
	"context"
)

type AiService interface {
	CreateModel(ctx context.Context, configuration model.CreateAIModel) *model.Ai
	GetModel(ctx context.Context) *model.Ai
	GenerateResponse(ctx context.Context, modelAI *model.Ai, histories *[]model.History, input string) string
}
