package controller

import "net/http"

type AuthController interface {
	Login(writer http.ResponseWriter, request *http.Request)
	Logout(writer http.ResponseWriter, request *http.Request)
}
