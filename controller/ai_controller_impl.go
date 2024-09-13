package controller

import (
	"adiubaidah/adi-bot/helper"
	"adiubaidah/adi-bot/model/ai"
	"adiubaidah/adi-bot/service"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

var activeRoutines = make(map[string]chan struct{})
var mu sync.Mutex // Protects access to the map

type AiControllerImpl struct {
	WaService     service.WaService
	AiService     service.AiService
	HistorService service.HistoryService
	Validate      *validator.Validate
}

func NewAiController(waService service.WaService, aiService service.AiService, historyService service.HistoryService, validate *validator.Validate) Controller {
	return &AiControllerImpl{
		WaService:     waService,
		AiService:     aiService,
		HistorService: historyService,
		Validate:      validate,
	}
}

func (a *AiControllerImpl) Activate(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	aiRequest := new(ai.AiRequest)
	helper.ReadFromRequestBody(request, aiRequest)
	err := a.Validate.Struct(aiRequest)
	helper.PanicIfError("Error validating request", err)

	// Goroutine to handle WhatsApp and AI service activation
	go func() {
		// Activate WhatsApp first
		waActive := a.WaService.Activate(request.Context(), aiRequest.Phone)

		// If WhatsApp is authenticated, proceed to activate AI
		if waActive {
			// Ensure only one instance per user
			mu.Lock()
			if _, exists := activeRoutines[aiRequest.Phone]; !exists {
				// Create AI Model
				// model := a.AiService.CreateModel(request.Context(), *aiRequest)

				// Start AI and WhatsApp in Go routine and store stop channel
				stopCh := make(chan struct{})
				activeRoutines[aiRequest.Phone] = stopCh
				go func() {
					for {
						select {
						case <-stopCh:
							// Handle stop signal (e.g., cleanup resources, stop AI service)
							fmt.Println("Stopping AI and WhatsApp service for", aiRequest.Phone)
							return
						default:

							histories, err := a.HistorService.GetHistory("6285232517546", aiRequest.Phone)
							helper.PanicIfError("Error getting history", err)
							fmt.Println(histories)

							// a.AiService.GenerateResponse(request.Context(), model, histories, "Hai Adi")
						}
					}
				}()
			}
			mu.Unlock()
		}
	}()

	helper.WriteToResponseBody(writer, "AI and WhatsApp are activating")
}

func (a *AiControllerImpl) Deactivate(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	phone := request.URL.Query().Get("phone")
	if phone == "" {
		helper.PanicIfError("Error deactivating AI and WhatsApp", fmt.Errorf("Phone number is required"))
		return
	}

	// Stop AI and WhatsApp service
	mu.Lock()
	if stopCh, exists := activeRoutines[phone]; exists {
		close(stopCh)
		delete(activeRoutines, phone)
	}
	mu.Unlock()

	a.WaService.Deactivate(request.Context(), phone)
	helper.WriteToResponseBody(writer, "AI and WhatsApp are deactivating")
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

func (a *AiControllerImpl) RegisterRoutes(router *httprouter.Router) {
	router.POST("/ai/activate", a.Activate)
	router.POST("/ai/deactivate", a.Deactivate)
	router.GET("/ai/check-activation", a.CheckActivation)
	router.GET("/ai/check-authentication", a.CheckAuthentication)
}
