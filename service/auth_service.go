package service

import (
	"context"

	"github.com/adiubaidah/wasabi/model"
)

type AuthService interface {
	Login(ctx context.Context, request model.UserLoginRequest) string
}
