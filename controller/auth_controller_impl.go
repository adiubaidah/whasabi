package controller

import (
	"net/http"
	"time"

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

	cookie := http.Cookie{
		Name:    "token",
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(time.Hour * 72),
	}
	http.SetCookie(writer, &cookie)
	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   "Login success",
	})
}

func (controller *AuthControllerImpl) IsAuth(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	userContext := request.Context().Value(middleware.UserContext).(jwt.MapClaims)
	userId := int(userContext["id"].(float64))
	if userId == 0 {
		helper.WriteToResponseBody(writer, &model.WebResponse{
			Code:   401,
			Status: "error",
			Data:   "Unauthorized",
		})
		return
	}

	user := controller.UserService.FindById(userId)

	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   user,
	})
}

func (controller *AuthControllerImpl) Logout(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	// Delete the token cookie
	cookie := http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",
	}

	http.SetCookie(writer, &cookie)

	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   "Logout success",
	})
}
