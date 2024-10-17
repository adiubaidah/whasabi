package controller

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strconv"

	"github.com/adiubaidah/wasabi/app"
	"github.com/adiubaidah/wasabi/exception"
	"github.com/adiubaidah/wasabi/helper"
	"github.com/adiubaidah/wasabi/middleware"
	"github.com/adiubaidah/wasabi/model"
	"github.com/adiubaidah/wasabi/service"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type ProcessControllerImpl struct {
	ProcessService service.ProcessService
	HistorService  service.HistoryService
	Validate       *validator.Validate
}

func NewProcessController(processService service.ProcessService, historyService service.HistoryService, validate *validator.Validate) ProcessController {
	return &ProcessControllerImpl{
		ProcessService: processService,
		HistorService:  historyService,
		Validate:       validate,
	}
}

func (a *ProcessControllerImpl) getSearchByUserId(request *http.Request, userContext jwt.MapClaims) (uint, error) {
	userRole := userContext["role"].(string)
	if userRole == "admin" {
		userIdStr := request.URL.Query().Get("user_id")
		userIdInt, err := strconv.Atoi(userIdStr)
		if err != nil {
			return 0, fmt.Errorf("error converting user id to int: %w", err)
		}
		return uint(userIdInt), nil
	}
	return uint(userContext["id"].(float64)), nil
}

func (a *ProcessControllerImpl) ListProcess(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	processes := a.ProcessService.ListProcess()
	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   processes,
	})
}

func (a *ProcessControllerImpl) GetModel(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	userContext := request.Context().Value(middleware.UserContext).(jwt.MapClaims)
	searchByUserId, err := a.getSearchByUserId(request, userContext)
	if err != nil {
		panic(exception.NewNotFoundError("Error find by user id"))
	}

	processModel := a.ProcessService.GetModel(searchByUserId)
	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   processModel,
	})
}

func (a *ProcessControllerImpl) UpsertModel(writter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	createProcessModel := new(model.CreateProcessModel)
	helper.ReadFromRequestBody(request, createProcessModel)
	err := a.Validate.Struct(createProcessModel)
	helper.PanicIfError("Error validating request", err)
	userContext := request.Context().Value(middleware.UserContext).(jwt.MapClaims)
	fmt.Println("Create Model", createProcessModel)
	searchByUserId, err := a.getSearchByUserId(request, userContext)
	if err != nil {
		panic(exception.NewNotFoundError("Error find by user id"))
	}
	// fmt.Println("searchByUserId", searchByUserId)

	result := a.ProcessService.UpsertModel(searchByUserId, *createProcessModel)

	helper.WriteToResponseBody(writter, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   result,
	})

}

func (a *ProcessControllerImpl) Activate(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	userContext := request.Context().Value(middleware.UserContext).(jwt.MapClaims)
	searchByUserId, err := a.getSearchByUserId(request, userContext)
	helper.PanicIfError("Error getting search by user id", err)

	//get user id from request context
	modelAi := a.ProcessService.GetModel(searchByUserId)

	if modelAi.Phone == "" {
		helper.WriteToResponseBody(writer, &model.WebResponse{
			Code:   400,
			Status: "error",
			Data:   "Phone number is not set",
		})
		return
	}
	fmt.Println("Active go routine after activation", runtime.NumGoroutine())

	go func() {

		app.Mu.Lock()
		defer app.Mu.Unlock()
		if _, exists := app.ActiveRoutines[modelAi.Phone]; !exists {
			waActive := a.ProcessService.Activate(modelAi.Phone)
			stopCh := make(chan struct{})
			app.ActiveRoutines[modelAi.Phone] = stopCh
			waActive.WaClient.AddEventHandler(func(evt any) {
				newContext := context.Background()
				switch v := evt.(type) {
				case *events.Message:
					if v.Info.IsGroup || v.Info.IsFromMe {
						return
					}

					if v.Message.GetConversation() == "" {
						return
					}

					if v.Info.Timestamp.After(waActive.StartTime) {
						fmt.Println("Pesan timestamp:", v.Info.Timestamp)
						fmt.Println("Pesan baru", v.Message.GetConversation())
						histories, err := a.HistorService.GetHistory(v.Info.Sender.String(), waActive.WaClient.Store.ID.String())
						helper.PanicIfError("Error getting history", err)
						input := v.Message.GetConversation()
						response := a.ProcessService.GenerateResponse(modelAi, histories, input)
						helper.PanicIfError("Error generating response", err)

						_, err = waActive.WaClient.SendMessage(newContext, v.Info.Chat, &waE2E.Message{
							Conversation: proto.String(response),
						})

						helper.PanicIfError("Error sending message", err)

						err = a.HistorService.InsertHistory(modelAi.ID, v.Info.Sender.String(), waActive.WaClient.Store.ID.String(), input, "user")
						helper.PanicIfError("Error inserting history user", err)
						err = a.HistorService.InsertHistory(modelAi.ID, waActive.WaClient.Store.ID.String(), v.Info.Sender.String(), response, "model")
						helper.PanicIfError("Error inserting history model", err)
					}
				default:
					fmt.Println("Event type not supported")
				}
			})

			go a.runAIService(stopCh, modelAi)
		}

	}()

	isAuthenticated := modelAi.IsAuthenticated
	isActivated := a.ProcessService.CheckActivation(modelAi.Phone)
	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data: map[string]any{
			"phone": modelAi.Phone,
			"status": map[string]bool{
				"is_authenticated": isAuthenticated,
				"is_activated":     isActivated,
			},
		},
	})

}

func (a *ProcessControllerImpl) runAIService(stopCh chan struct{}, model *model.Process) {

	for { // it will keep running until the stop channel is closed
		select {
		case <-stopCh:
			fmt.Println("Stopping Goroutine for phone:", model.Phone)

			return
		}
	}
}

func (a *ProcessControllerImpl) Deactivate(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userContext := request.Context().Value(middleware.UserContext).(jwt.MapClaims)
	searchByUserId, err := a.getSearchByUserId(request, userContext)
	helper.PanicIfError("Error getting search by user id", err)

	modelProcess := a.ProcessService.GetModel(searchByUserId)
	a.ProcessService.Deactivate(modelProcess.Phone)

	app.Mu.Lock()
	defer app.Mu.Unlock() // Ensure we unlock the mutex at the end

	if stopCh, exists := app.ActiveRoutines[modelProcess.Phone]; exists {
		a.ProcessService.Deactivate(modelProcess.Phone)
		close(stopCh)                                  // Close the stop channel to signal the Goroutine to stop
		delete(app.ActiveRoutines, modelProcess.Phone) // Remove the phone from the map
	} else {
		helper.WriteToResponseBody(writer, &model.WebResponse{
			Code:   400,
			Status: "error",
			Data:   "AI service is not active",
		})
		return
	}

	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   "AI and WhatsApp are deactivating",
	})
	fmt.Println("Number of Goroutines after deactivation:", runtime.NumGoroutine())
}

func (a *ProcessControllerImpl) CheckActivation(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userContext := request.Context().Value(middleware.UserContext).(jwt.MapClaims)
	searchByUserId, err := a.getSearchByUserId(request, userContext)
	if err != nil {
		panic(exception.NewNotFoundError("Error find by user id"))
	}

	modelProcess := a.ProcessService.GetModel(searchByUserId)
	status := a.ProcessService.CheckActivation(modelProcess.Phone)
	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   status,
	})

}

func (a *ProcessControllerImpl) CheckAuthentication(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userContext := request.Context().Value(middleware.UserContext).(jwt.MapClaims)
	searchByUserId, err := a.getSearchByUserId(request, userContext)
	if err != nil {
		panic(exception.NewNotFoundError("Error find by user id"))
	}

	modelProcess := a.ProcessService.GetModel(searchByUserId)
	status := a.ProcessService.CheckAuthentication(modelProcess.Phone)
	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   status,
	})
}
func (a *ProcessControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userContext := request.Context().Value(middleware.UserContext).(jwt.MapClaims)
	searchByUserId, err := a.getSearchByUserId(request, userContext)
	if err != nil {
		panic(exception.NewNotFoundError("Error find by user id"))
	}

	modelProcess := a.ProcessService.GetModel(searchByUserId)
	if a.ProcessService.CheckActivation(modelProcess.Phone) {
		helper.WriteToResponseBody(writer, &model.WebResponse{
			Code:   400,
			Status: "error",
			Data:   "Please deactivate the AI service first",
		})
		return
	}
	result := a.ProcessService.Delete(modelProcess.Phone)
	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   result,
	})
}
