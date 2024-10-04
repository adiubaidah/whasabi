package service

import (
	"context"

	"github.com/adiubaidah/wasabi/helper"
	"github.com/adiubaidah/wasabi/model"

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
		Username: request.Username,
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

func (service *UserServiceImpl) FindById(userId int) *model.UserDTO {
	var user model.UserDTO
	err := service.DB.Table("users").Select("id", "username", "role").Where("id = ?", userId).First(&user).Error
	helper.PanicIfError("Error finding user by ID", err)
	return &user
}

func (service *UserServiceImpl) FindByUsername(username string) *model.UserDTO {
	var user model.UserDTO
	err := service.DB.Table("users").Select("id", "username", "role").Where("username = ?", username).First(&user).Error
	helper.PanicIfError("Error finding user by username", err)
	return &user
}

func (service *UserServiceImpl) Find(params UserSearchParams) *[]model.UserDTO {
	var users []model.UserDTO
	query := service.DB.Table("users").Select("id", "username", "role")

	if params.Role != "" {
		query = query.Where("role = ?", params.Role)
	}

	err := query.Find(&users).Error
	helper.PanicIfError("Error finding users", err)

	return &users
}
