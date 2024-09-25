package controller

import (
	"adiubaidah/adi-bot/helper"
	"adiubaidah/adi-bot/model"
	"adiubaidah/adi-bot/service"
	"context"
	"fmt"
	"net/http"
	"runtime"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/google/generative-ai-go/genai"
	"github.com/julienschmidt/httprouter"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

var activeRoutines = make(map[string]chan struct{})
var mu sync.Mutex // Protects access to the map

type AiControllerImpl struct {
	WaService     service.WaService
	AiService     service.AiService
	HistorService service.HistoryService
	Validate      *validator.Validate
	AiModel       *genai.GenerativeModel
}

func NewAiController(waService service.WaService, aiService service.AiService, historyService service.HistoryService, validate *validator.Validate) AiController {
	return &AiControllerImpl{
		WaService:     waService,
		AiService:     aiService,
		HistorService: historyService,
		Validate:      validate,
	}
}

func (a *AiControllerImpl) GetModel(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ai := a.AiService.GetModel(request.Context())
	helper.WriteToResponseBody(writer, ai)
}

func (a *AiControllerImpl) CreateModel(writter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	createModelAi := new(model.CreateAIModel)
	helper.ReadFromRequestBody(request, createModelAi)
	err := a.Validate.Struct(createModelAi)
	helper.PanicIfError("Error validating request", err)

	result := a.AiService.CreateModel(request.Context(), *createModelAi)

	helper.WriteToResponseBody(writter, result)

}

func (a *AiControllerImpl) Activate(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	//get user id from request context
	modelAi := a.AiService.GetModel(request.Context())

	go func() {

		mu.Lock()
		defer mu.Unlock()
		if _, exists := activeRoutines[modelAi.Phone]; !exists {
			waActive := a.WaService.Activate(request.Context(), modelAi.Phone)
			// Start a new routine only if it doesn't already exist
			stopCh := make(chan struct{})
			activeRoutines[modelAi.Phone] = stopCh
			waActive.WaClient.AddEventHandler(func(evt any) {
				newContext := context.Background()
				switch v := evt.(type) {
				case *events.Message:
					if v.Info.IsGroup {
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
						response := a.AiService.GenerateResponse(newContext, modelAi, histories, input)
						helper.PanicIfError("Error generating response", err)

						_, err = waActive.WaClient.SendMessage(newContext, v.Info.Chat, &waE2E.Message{
							Conversation: proto.String(response),
						})

						helper.PanicIfError("Error sending message", err)

						err = a.HistorService.InsertHistory(v.Info.Sender.String(), waActive.WaClient.Store.ID.String(), input, "user")
						helper.PanicIfError("Error inserting history user", err)
						err = a.HistorService.InsertHistory(waActive.WaClient.Store.ID.String(), v.Info.Sender.String(), response, "model")
						helper.PanicIfError("Error inserting history model", err)
					}
				default:
					fmt.Println("Event type not supported")
				}
			})

			go a.runAIService(stopCh, modelAi)
		}
	}()

	helper.WriteToResponseBody(writer, "AI and WhatsApp are activating")

}

func (a *AiControllerImpl) runAIService(stopCh chan struct{}, model *model.Ai) {

	for { // it will keep running until the stop channel is closed
		select {
		case <-stopCh:
			fmt.Println("Stopping Goroutine for phone:", model.Phone)

			return
		}
	}
}

func (a *AiControllerImpl) Deactivate(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	phone := request.URL.Query().Get("phone")
	if phone == "" {
		helper.PanicIfError("Error deactivating AI and WhatsApp", fmt.Errorf("Phone number is required"))
		return
	}

	mu.Lock()
	defer mu.Unlock() // Ensure we unlock the mutex at the end

	if stopCh, exists := activeRoutines[phone]; exists {
		a.WaService.Deactivate(request.Context(), phone)
		close(stopCh)                 // Close the stop channel to signal the Goroutine to stop
		delete(activeRoutines, phone) // Remove the phone from the map
	} else {
		helper.WriteToResponseBody(writer, "No active session found for this phone")
		return
	}

	helper.WriteToResponseBody(writer, "AI and WhatsApp are deactivating")
	fmt.Println("Number of Goroutines after deactivation:", runtime.NumGoroutine())
}

func (a *AiControllerImpl) CheckActivation(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	phone := request.URL.Query().Get("phone")
	if phone == "" {
		helper.PanicIfError("Error checking activation", fmt.Errorf("Phone number is required"))
		return
	}
	status := a.WaService.CheckActivation(request.Context(), phone)
	helper.WriteToResponseBody(writer, fmt.Sprintf("AI and WhatsApp are active: %t", status))

}

func (a *AiControllerImpl) CheckAuthentication(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	phone := request.URL.Query().Get("phone")
	if phone == "" {
		helper.PanicIfError("Error checking authentication", fmt.Errorf("Phone number is required"))
		return
	}
	status := a.WaService.CheckAuthentication(request.Context(), phone)
	helper.WriteToResponseBody(writer, fmt.Sprintf("AI and WhatsApp are authenticated: %t", status))
}
