package service

import (
	"context"

	"github.com/adiubaidah/wasabi/model"
)

type UserService interface {
	Create(ctx context.Context, request model.UserCreateRequest) *model.User
	FindById(userId int) *model.UserDTO
	FindByUsername(username string) *model.UserDTO
	Find(params UserSearchParams) *[]model.UserDTO
}
