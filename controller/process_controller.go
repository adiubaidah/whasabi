package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ProcessController interface {
	ListProcess(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetModel(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpsertModel(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Activate(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Deactivate(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)

	CheckActivation(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CheckAuthentication(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
