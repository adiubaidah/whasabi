package service

import (
	"adiubaidah/adi-bot/exception"
	"adiubaidah/adi-bot/helper"
	"adiubaidah/adi-bot/model"
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthServiceImpl struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthServiceImpl {
	return &AuthServiceImpl{
		DB: db,
	}
}

func (a *AuthServiceImpl) Login(ctx context.Context, request model.UserLoginRequest) string {
	// Validate the request

	// Check if the user exists
	user := model.User{}
	a.DB.Where("username = ?", request.Username).Take(&user)
	// Compare the password with the hashed password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
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

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	helper.PanicIfError("Error signing the token", err)
	return tokenString
}
