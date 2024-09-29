package service

import (
	"adiubaidah/adi-bot/model"
	"context"
)

type AuthService interface {
	Login(ctx context.Context, request model.UserLoginRequest) string
}
