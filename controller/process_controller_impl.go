package controller

import (
	"adiubaidah/adi-bot/app"
	"adiubaidah/adi-bot/helper"
	"adiubaidah/adi-bot/middleware"
	"adiubaidah/adi-bot/model"
	"adiubaidah/adi-bot/service"
	"context"
	"fmt"
	"net/http"
	"runtime"

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

func (a *ProcessControllerImpl) GetModel(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	userContext := request.Context().Value(middleware.UserContext).(jwt.MapClaims)
	userId := uint(userContext["id"].(float64))

	processModel := a.ProcessService.GetModel(userId)
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
	userId := uint(userContext["id"].(float64))

	result := a.ProcessService.UpsertModel(userId, *createProcessModel)

	helper.WriteToResponseBody(writter, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   result,
	})

}

func (a *ProcessControllerImpl) Activate(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	userContext := request.Context().Value(middleware.UserContext).(jwt.MapClaims)
	userId := uint(userContext["id"].(float64))

	//get user id from request context
	modelAi := a.ProcessService.GetModel(userId)
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

	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   "AI and WhatsApp are activating",
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
	userId := uint(userContext["id"].(float64))

	modelProcess := a.ProcessService.GetModel(userId)
	a.ProcessService.Deactivate(modelProcess.Phone)

	app.Mu.Lock()
	defer app.Mu.Unlock() // Ensure we unlock the mutex at the end

	if stopCh, exists := app.ActiveRoutines[modelProcess.Phone]; exists {
		a.ProcessService.Deactivate(modelProcess.Phone)
		close(stopCh)                                  // Close the stop channel to signal the Goroutine to stop
		delete(app.ActiveRoutines, modelProcess.Phone) // Remove the phone from the map
	} else {
		helper.WriteToResponseBody(writer, "No active session found for this phone")
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
	userId := uint(userContext["id"].(float64))

	modelProcess := a.ProcessService.GetModel(userId)
	status := a.ProcessService.CheckActivation(modelProcess.Phone)
	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   status,
	})

}

func (a *ProcessControllerImpl) CheckAuthentication(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userContext := request.Context().Value(middleware.UserContext).(jwt.MapClaims)
	userId := uint(userContext["id"].(float64))

	modelProcess := a.ProcessService.GetModel(userId)
	status := a.ProcessService.CheckAuthentication(modelProcess.Phone)
	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   status,
	})
}
