package controller

import "github.com/julienschmidt/httprouter"

type Controller interface {
	RegisterRoutes(router *httprouter.Router)
}
