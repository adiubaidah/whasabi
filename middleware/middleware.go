package middleware

import (
	"context"
	"net/http"

	"github.com/adiubaidah/wasabi/exception"
	"github.com/adiubaidah/wasabi/helper"

	"github.com/golang-jwt/jwt/v5"
)

type Middleware struct {
	http.Handler
}

type contextKey string

const UserContext contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			panic(exception.NewUnauthorizedError("Unauthorized"))
		}

		tokenString := cookie.Value
		token, err := helper.JwtParse(tokenString)

		if err != nil || !token.Valid {
			panic(exception.NewUnauthorizedError("Token tidak valid"))
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			panic(exception.NewUnauthorizedError("Token tidak valid"))
		}

		ctx := context.WithValue(r.Context(), UserContext, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userContext := r.Context().Value(UserContext).(jwt.MapClaims)

		if userContext["role"] != "admin" {
			panic(exception.NewForbiddenError("Forbidden"))
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Handler.ServeHTTP(w, r)
}
