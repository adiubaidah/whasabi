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

func SetupRouter(processCtrl controller.ProcessController, userCtrl controller.UserController, authCtrl controller.AuthController) *httprouter.Router {
	router := httprouter.New()
	router.Handler(http.MethodPost, "/auth/login", adaptHandle(authCtrl.Login))
	router.Handler(http.MethodGet, "/auth/is-auth", middleware.AuthMiddleware(adaptHandle(authCtrl.IsAuth)))
	router.Handler(http.MethodPost, "/auth/logout", middleware.AuthMiddleware(adaptHandle(authCtrl.Logout)))

	// Wrap AI controller routes with AuthMiddleware (for authenticated users)
	router.Handler(http.MethodGet, "/process/model", middleware.AuthMiddleware(adaptHandle(processCtrl.GetModel)))
	router.Handler(http.MethodPost, "/process/model", middleware.AuthMiddleware(adaptHandle(processCtrl.UpsertModel)))
	router.Handler(http.MethodPost, "/process/activate", middleware.AuthMiddleware(adaptHandle(processCtrl.Activate)))
	router.Handler(http.MethodPost, "/process/deactivate", middleware.AuthMiddleware(adaptHandle(processCtrl.Deactivate)))
	router.Handler(http.MethodGet, "/process/check-activation", middleware.AuthMiddleware(adaptHandle(processCtrl.CheckActivation)))
	router.Handler(http.MethodGet, "/process/check-authentication", middleware.AuthMiddleware(adaptHandle(processCtrl.CheckAuthentication)))

	// Wrap User controller routes with AdminMiddleware (for admin role only)
	router.Handler(http.MethodPost, "/user", middleware.AuthMiddleware(middleware.AdminMiddleware(adaptHandle(userCtrl.Create))))
	router.Handler(http.MethodGet, "/user-ws", adaptHandle(userCtrl.WebSocket))
	router.ServeFiles("/public/*filepath", http.Dir("public"))

	router.PanicHandler = exception.ErrorHandler
	return router
}
