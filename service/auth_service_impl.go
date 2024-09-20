package service

import (
	"adiubaidah/adi-bot/helper"
	"adiubaidah/adi-bot/model/user"
	"context"
	"database/sql"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type AuthServiceImpl struct {
	DB       *sql.DB
	Validate *validator.Validate
}

func NewAuthService(db *sql.DB, validate *validator.Validate) *AuthServiceImpl {
	return &AuthServiceImpl{
		DB:       db,
		Validate: validate,
	}
}

func (service *AuthServiceImpl) Login(ctx context.Context, request user.UserLoginRequest) (string, error) {
	// Validate the request
	err := service.Validate.Struct(request)
	if err != nil {
		return "", err
	}
	var userID int
	var role string
	err = service.DB.QueryRowContext(ctx, "SELECT id, role FROM users WHERE username = ? AND password = ?", request.Username, request.Password).Scan(&userID, &role)
	if err != nil {
		return "", err
	}

	// Create JWT claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// return token, nil

	tokenString, err := token.SignedString([]byte(helper.GetEnv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (service *AuthServiceImpl) Logout(ctx context.Context, request user.UserLogoutRequest) error {
	// Validate the request
	err := service.Validate.Struct(request)
	if err != nil {
		return err
	}

	//

	// return nil
	return nil
}
