package exception

import (
	"log"
	"net/http"

	"github.com/adiubaidah/wasabi/helper"
	"github.com/adiubaidah/wasabi/model"

	"github.com/go-playground/validator/v10"
)

func ErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {

	log.Default().Println(err)

	if notFoundError(writer, request, err) {
		return
	}

	if badRequestError(writer, request, err) {
		return
	}

	if validationErrors(writer, request, err) {
		return
	}

	if notAuthorizedError(writer, request, err) {
		return
	}

	if forbiddenError(writer, request, err) {
		return
	}

	internalServerError(writer, request, err)

	// Handle unhandled errors (including generic panics)

}

func validationErrors(writer http.ResponseWriter, _ *http.Request, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors) //if convertion is success, ok will be true
	if ok {
		writer.Header().Set("Content-Type", "application/json")

		webResponse := &model.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   exception.Error(),
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	}
	return false
}

func badRequestError(writer http.ResponseWriter, _ *http.Request, err interface{}) bool {
	exception, ok := err.(BadRequestError)
	if ok {
		writer.Header().Set("Content-Type", "application/json")

		webResponse := &model.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   exception.Error,
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	}
	return false
}

func notFoundError(writer http.ResponseWriter, _ *http.Request, err interface{}) bool {
	exception, ok := err.(NotFoundError)
	if ok {
		writer.Header().Set("Content-Type", "application/json")

		webResponse := &model.WebResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Data:   exception.Error,
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	}
	return false
}

func notAuthorizedError(writer http.ResponseWriter, _ *http.Request, err interface{}) bool {
	exception, ok := err.(UnauthorizedError)
	if ok {
		writer.Header().Set("Content-Type", "application/json")

		webResponse := &model.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
			Data:   exception.Error,
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	}
	return false
}

func forbiddenError(writer http.ResponseWriter, _ *http.Request, err interface{}) bool {
	exception, ok := err.(ForbiddenError)
	if ok {
		writer.Header().Set("Content-Type", "application/json")

		webResponse := &model.WebResponse{
			Code:   http.StatusForbidden,
			Status: "FORBIDDEN",
			Data:   exception.Error,
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	}
	return false
}

func internalServerError(writer http.ResponseWriter, _ *http.Request, err interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	webResponse := &model.WebResponse{
		Code:   http.StatusInternalServerError,
		Status: "INTERNAL SERVER ERROR",
		Data:   err,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
