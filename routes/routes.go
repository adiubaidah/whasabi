package routes

import (
	"adiubaidah/adi-bot/controller"
	"adiubaidah/adi-bot/exception"

	"github.com/julienschmidt/httprouter"
)

func SetupRouter(controllers ...controller.Controller) *httprouter.Router {
	router := httprouter.New()

	for _, ctrl := range controllers {
		ctrl.RegisterRoutes(router)
	}

	router.PanicHandler = exception.ErrorHandler

	return router
}
