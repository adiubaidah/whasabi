package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/adiubaidah/wasabi/exception"
	"github.com/adiubaidah/wasabi/helper"
	"github.com/adiubaidah/wasabi/middleware"
	"github.com/adiubaidah/wasabi/model"
	"github.com/adiubaidah/wasabi/service"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

type AuthControllerImpl struct {
	AuthService service.AuthService
	service.UserService
	Validate *validator.Validate
}

func NewAuthController(authService service.AuthService, userService service.UserService, validate *validator.Validate) AuthController {
	return &AuthControllerImpl{
		AuthService: authService,
		UserService: userService,
		Validate:    validate,
	}
}

func (controller *AuthControllerImpl) Login(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userLoginRequest := new(model.UserLoginRequest)
	helper.ReadFromRequestBody(request, userLoginRequest)

	err := controller.Validate.Struct(userLoginRequest)
	helper.PanicIfError("Error validating request", err)
	token := controller.AuthService.Login(request.Context(), *userLoginRequest)
	cookieSameSite := http.SameSiteDefaultMode

	if os.Getenv("ENVIROMENT") == "production" {
		cookieSameSite = http.SameSiteNoneMode
	}

	cookie := http.Cookie{
		Name:     os.Getenv("COOKIE_NAME"),
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 72),
		SameSite: cookieSameSite,
		Domain:   os.Getenv("COOKIE_DOMAIN"),
	}
	http.SetCookie(writer, &cookie)
	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   "Login success",
	})
}

func (controller *AuthControllerImpl) Register(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userRegister := new(model.UserRegisterRequest)
	helper.ReadFromRequestBody(request, userRegister)

	err := controller.Validate.Struct(userRegister)
	helper.PanicIfError("Error validating request", err)

	user := controller.UserService.Create(request.Context(), &model.UserCreateRequest{
		Username: userRegister.Username,
		Password: userRegister.Password,
		Role:     "user",
		IsActive: false,
	})

	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data: &model.UserDTO{
			ID:       user.ID,
			Username: user.Username,
			Role:     user.Role,
			IsActive: user.IsActive,
		},
	})
}

func (controller *AuthControllerImpl) IsAuth(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	userContext := request.Context().Value(middleware.UserContext).(jwt.MapClaims)
	userId := int(userContext["id"].(float64))
	if userId == 0 {
		panic(exception.NewUnauthorizedError("Unauthorized"))
	}

	user := controller.UserService.FindById(userId)

	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   user,
	})
}

func (controller *AuthControllerImpl) Logout(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	cookieSameSite := http.SameSiteDefaultMode
	if os.Getenv("ENVIROMENT") == "production" {
		cookieSameSite = http.SameSiteNoneMode
	}
	// Delete the token cookie
	cookie := http.Cookie{
		Name:     os.Getenv("COOKIE_NAME"),
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		SameSite: cookieSameSite,
		Domain:   os.Getenv("COOKIE_DOMAIN"),
	}

	http.SetCookie(writer, &cookie)

	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   "Logout success",
	})
}
