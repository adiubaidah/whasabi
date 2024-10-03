package service

import (
	"adiubaidah/adi-bot/helper"
	"adiubaidah/adi-bot/model"
	"context"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &UserServiceImpl{
		DB: db,
	}
}

// Create is a function to create a new user
func (service *UserServiceImpl) Create(ctx context.Context, request model.UserCreateRequest) *model.User {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	helper.PanicIfError("Erro generating password", err)

	user := model.User{
		Username: request.Password,
		Password: string(hashedPassword),
		Role:     request.Role,
	}

	err = service.DB.Create(&user).Error
	helper.PanicIfError("Error create user", err)

	return &user

}

type UserSearchParams struct {
	UserId   int
	Username string
	Role     string
}

func (service *UserServiceImpl) Find(params UserSearchParams) *model.UserDTO {
	user := model.UserDTO{}
	query := service.DB.Table("users").Select("id", "username", "role")
	if params.UserId != 0 {
		query = query.Where("id = ?", params.UserId)
	}
	if params.Username != "" {
		query = query.Where("username = ?", params.Username)
	}
	if params.Role != "" {
		query = query.Where("role = ?", params.Role)
	}

	err := query.Take(&user).Error
	helper.PanicIfError("Error finding user", err)

	return &user
}
