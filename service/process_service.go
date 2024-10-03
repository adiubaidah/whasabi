package service

import (
	"adiubaidah/adi-bot/model"
)

type ProcessService interface {
	Activate(phone string) *UserWaStatus
	CheckActivation(phone string) bool
	CheckAuthentication(phone string) bool
	Deactivate(phone string) bool

	UpsertModel(userId uint, configuration model.CreateProcessModel) *model.Process
	GetModel(userId uint) *model.Process
	GenerateResponse(modelAI *model.Process, histories *[]model.History, input string) string
}
