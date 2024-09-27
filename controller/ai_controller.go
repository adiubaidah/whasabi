package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type AiController interface {
	GetModel(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CreateModel(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Activate(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Deactivate(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CheckActivation(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CheckAuthentication(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
