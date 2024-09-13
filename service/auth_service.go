package service

import (
	"adiubaidah/adi-bot/model/user"
	"context"
)

type AuthService interface {
	Login(ctx context.Context, request user.UserLoginRequest) (string, error)
	Logout(ctx context.Context, request user.UserLogoutRequest) error
}
