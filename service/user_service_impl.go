package service

import (
	"adiubaidah/adi-bot/helper"
	"adiubaidah/adi-bot/model"
	"context"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewUserService(db *gorm.DB, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		DB:       db,
		Validate: validate,
	}
}

// Create is a function to create a new user
func (service *UserServiceImpl) Create(ctx context.Context, request model.UserCreateRequest) model.User {
	err := service.Validate.Struct(request)
	helper.PanicIfError("", err)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	helper.PanicIfError("", err)

	user := model.User{
		Username: request.Username,
		Password: string(hashedPassword),
		Role:     request.Role,
	}

	err = service.DB.Create(&user).Error
	helper.PanicIfError("", err)

	return user
}
