package service

import (
	"adiubaidah/adi-bot/model"
	"context"
)

type UserService interface {
	Create(ctx context.Context, request model.UserCreateRequest) *model.User
	Find(params UserSearchParams) *model.UserDTO
}
