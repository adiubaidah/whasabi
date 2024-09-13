package service

import (
	"adiubaidah/adi-bot/model/user"
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
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
	return "", nil
	// return token, nil
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
