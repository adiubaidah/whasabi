package service

import (
	"adiubaidah/adi-bot/helper"
	"adiubaidah/adi-bot/model"
	"context"
	"fmt"

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
func (service *UserServiceImpl) Create(ctx context.Context, request model.UserCreateRequest) *model.User {
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
	fmt.Println(params)
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

func (service *UserServiceImpl) GetService(userId int) *model.Ai {

	ai := model.Ai{}

	err := service.DB.Where("user_id = ?", userId).Take(&ai).Error
	helper.PanicIfError("Error getting service by user id", err)

	return &ai
}
