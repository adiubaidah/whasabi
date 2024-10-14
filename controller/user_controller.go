package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type UserController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writter http.ResponseWriter, request *http.Request, params httprouter.Params)
	List(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	WebSocket(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
