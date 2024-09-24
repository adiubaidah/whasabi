package service

import (
	"adiubaidah/adi-bot/exception"
	"adiubaidah/adi-bot/helper"
	"adiubaidah/adi-bot/model"
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthServiceImpl struct {
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewAuthService(db *gorm.DB, validate *validator.Validate) *AuthServiceImpl {
	return &AuthServiceImpl{
		DB:       db,
		Validate: validate,
	}
}

func (a *AuthServiceImpl) Login(ctx context.Context, request model.UserLoginRequest) string {
	// Validate the request
	err := a.Validate.Struct(request)
	helper.PanicIfError("", err)

	// Check if the user exists
	user := model.User{}
	a.DB.Where("username = ?", request.Username).Take(&user)
	// Compare the password with the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		panic(exception.NewUnauthorizedError("Invalid username or password"))
	}

	// Create JWT claims
	claims := jwt.MapClaims{
		"id":   user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(helper.GetEnv("JWT_SECRET")))
	helper.PanicIfError("Error signing the token", err)
	return tokenString
}

func (service *AuthServiceImpl) Logout(ctx context.Context, request model.UserLogoutRequest) error {
	// Validate the request
	err := service.Validate.Struct(request)
	if err != nil {
		return err
	}

	//

	// return nil
	return nil
}
