package routes

import (
	"adiubaidah/adi-bot/controller"
	"adiubaidah/adi-bot/exception"
	"adiubaidah/adi-bot/middleware"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Adapter function to convert httprouter.Handle to http.HandlerFunc because middleware.AuthMiddleware expects http.HandlerFunc
func adaptHandle(h httprouter.Handle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		// user := r.Context().Value(middleware.UserContext)
		h(w, r, params)
	}
}

func SetupRouter(aiCtrl controller.AiController, userCtrl controller.UserController, authCtrl controller.AuthController) *httprouter.Router {
	router := httprouter.New()
	router.Handler(http.MethodPost, "/login", adaptHandle(authCtrl.Login))

	// Wrap AI controller routes with AuthMiddleware (for authenticated users)
	router.Handler(http.MethodGet, "/ai/model", middleware.AuthMiddleware(adaptHandle(aiCtrl.GetModel)))
	router.Handler(http.MethodPost, "/ai/model", middleware.AuthMiddleware(adaptHandle(aiCtrl.CreateModel)))
	router.Handler(http.MethodPost, "/ai/activate", middleware.AuthMiddleware(adaptHandle(aiCtrl.Activate)))
	router.Handler(http.MethodPost, "/ai/deactivate", middleware.AuthMiddleware(adaptHandle(aiCtrl.Deactivate)))
	router.Handler(http.MethodGet, "/ai/check-activation", middleware.AuthMiddleware(adaptHandle(aiCtrl.CheckActivation)))
	router.Handler(http.MethodGet, "/ai/check-authentication", middleware.AuthMiddleware(adaptHandle(aiCtrl.CheckAuthentication)))

	// Wrap User controller routes with AdminMiddleware (for admin role only)
	router.Handler(http.MethodPost, "/user", middleware.AuthMiddleware(middleware.AdminMiddleware(adaptHandle(userCtrl.Create))))
	router.Handler(http.MethodGet, "/user-ws", adaptHandle(userCtrl.WebSocket))

	router.PanicHandler = exception.ErrorHandler
	return router
}
