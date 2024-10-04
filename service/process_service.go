package service

import (
	"github.com/adiubaidah/wasabi/model"
)

type ProcessService interface {
	ListProcess() *[]model.ProcessWithUserDTO
	Activate(phone string) *UserWaStatus
	CheckActivation(phone string) bool
	CheckAuthentication(phone string) bool
	Deactivate(phone string) bool
	Delete(phone string) bool

	UpsertModel(userId uint, configuration model.CreateProcessModel) *model.Process
	GetModel(userId uint) *model.Process
	GenerateResponse(modelAI *model.Process, histories *[]model.History, input string) string
}
