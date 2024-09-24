package controller

import (
	"adiubaidah/adi-bot/helper"
	"adiubaidah/adi-bot/model"
	"adiubaidah/adi-bot/service"
	"fmt"
	"net/http"
	"runtime"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/google/generative-ai-go/genai"
	"github.com/julienschmidt/httprouter"
	"go.mau.fi/whatsmeow/types/events"
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

func (a *AiControllerImpl) GetConfiguration(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ai := a.AiService.GetConfiguration(request.Context())
	helper.WriteToResponseBody(writer, ai)
}

func (a *AiControllerImpl) CreateConfiguration(writter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	aiConfiguration := new(model.AiConfiguration)
	helper.ReadFromRequestBody(request, aiConfiguration)
	err := a.Validate.Struct(aiConfiguration)
	helper.PanicIfError("Error validating request", err)

	result := a.AiService.CreateConfiguration(request.Context(), *aiConfiguration)

	helper.WriteToResponseBody(writter, result)

}

func (a *AiControllerImpl) Activate(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	aiConfiguration := new(model.AiConfiguration)
	helper.ReadFromRequestBody(request, aiConfiguration)
	err := a.Validate.Struct(aiConfiguration)
	helper.PanicIfError("Error validating request", err)
	// modela.AiService.CreateModel(request.Context(), *aiConfiguration)
	go func() {

		mu.Lock()
		defer mu.Unlock()
		if _, exists := activeRoutines[aiConfiguration.Phone]; !exists {
			waActive := a.WaService.Activate(request.Context(), aiConfiguration.Phone)
			// Start a new routine only if it doesn't already exist
			stopCh := make(chan struct{})
			activeRoutines[aiConfiguration.Phone] = stopCh

			go a.runAIService(stopCh, waActive, aiConfiguration.Phone)
		}
	}()

	helper.WriteToResponseBody(writer, "AI and WhatsApp are activating")

}

func (a *AiControllerImpl) runAIService(stopCh chan struct{}, waActive *service.UserWaStatus, phone string) {
	waActive.WaClient.AddEventHandler(func(evt any) {
		switch v := evt.(type) {
		case *events.Message:
			if v.Info.IsGroup {
				return
			}

			if v.Info.Timestamp.After(waActive.StartTime) {
				fmt.Println("Pesan timestamp:", v.Info.Timestamp)
				fmt.Println("Pesan baru", v.Message.GetConversation())
				// context := context.Background()
				// histories, err := a.HistorService.GetHistory(v.Info.Sender.String(), waActive.WaClient.Store.ID.String())
				// helper.PanicIfError("Error getting history", err)
				// input := v.Message.GetConversation()
				// response, err := a.AiService.GenerateResponse(context, a., histories, input)
			}
		}
	})

	for { // it will keep running until the stop channel is closed
		select {
		case <-stopCh:
			fmt.Println("Stopping Goroutine for phone:", phone)

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
