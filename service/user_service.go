package service

import (
	"adiubaidah/adi-bot/model"
	"context"
)

type UserService interface {
	Create(ctx context.Context, request model.UserCreateRequest) *model.User
	GetService(userId int) *model.Ai
}
