package service

import (
	"adiubaidah/adi-bot/app"
	"adiubaidah/adi-bot/helper"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

// Struct for storing user service status
type UserWaStatus struct {
	WaClient        *whatsmeow.Client
	Container       *sqlstore.Container
	IsActive        bool
	IsAuthenticated bool
	StartTime       time.Time
}

type WaServiceImpl struct {
	UserStatusMap map[string]*UserWaStatus // Map to track status by phone number
	WebSocketHub  *app.WebSocketHub
	mu            sync.Mutex // Mutex to protect concurrent access to map
}

// Create new WhatsApp and AI service
func NewWaService(waWebSocketHub *app.WebSocketHub) WaService {
	return &WaServiceImpl{
		UserStatusMap: make(map[string]*UserWaStatus),
		WebSocketHub:  waWebSocketHub,
	}
}

// Function to activate WhatsApp and AI for a user
func (s *WaServiceImpl) Activate(phone string) *UserWaStatus {
	// If the service is already active, no need to activate again
	s.mu.Lock()
	defer s.mu.Unlock()
	if status, exists := s.UserStatusMap[phone]; exists && status.IsActive {
		return status
	}

	waClient, container := app.GetWaClient(phone)
	// Initialize a new status if it doesn't exist
	s.UserStatusMap[phone] = &UserWaStatus{
		WaClient:        waClient,
		Container:       container,
		IsActive:        false,
		IsAuthenticated: false,
		StartTime:       time.Now(),
	}

	// Set the current status variable for convenience
	status := s.UserStatusMap[phone]
	status.IsAuthenticated = status.WaClient.Store.ID != nil
	status.IsActive = true
	fmt.Println("Wa Status", status.IsAuthenticated)

	// If not authenticated, handle the authentication process
	if !status.IsAuthenticated {
		s.handleQRCodeAuthentication(phone, status)
	} else {
		status.WaClient.Connect()
		s.WebSocketHub.SendMessage(phone, "Connection Successful!")
	}

	return status
}

// handleQRCodeAuthentication manages the WhatsApp authentication process using QR code
func (s *WaServiceImpl) handleQRCodeAuthentication(phone string, status *UserWaStatus) {
	context := context.Background()
	qrChan, err := status.WaClient.GetQRChannel(context)
	helper.PanicIfError("Error getting QR channel", err)
	err = status.WaClient.Connect()
	helper.PanicIfError("Error connecting to WhatsApp", err)
	go func() {
		for evt := range qrChan {
			fmt.Println("QR event received:", evt.Event)
			switch evt.Event {
			case "code":
				qrPath := fmt.Sprintf("public/qr-%s-%s.png", time.Now().Format("20060102-150405"), phone)
				err := qrcode.WriteFile(evt.Code, qrcode.Medium, 256, qrPath)
				if err != nil {
					fmt.Println("Error generating QR code:", err)
				} else {
					fmt.Printf("QR code generated for phone %s. Scan it using WhatsApp!\n", phone)
					s.WebSocketHub.SendMessage(phone, qrPath)
				}
			case "success":
				s.mu.Lock()
				defer s.mu.Unlock()
				status.IsAuthenticated = true
				s.WebSocketHub.SendMessage(phone, "Connection Successful!")
				return
			case "timeout":
				fmt.Println("Timeout")
				status.WaClient.Disconnect()
				s.WebSocketHub.SendMessage(phone, "Connection Timeout!")
				if stopCh, exists := app.ActiveRoutines[phone]; exists {
					s.Deactivate(context, phone)
					close(stopCh)                     // Close the stop channel to signal the Goroutine to stop
					delete(app.ActiveRoutines, phone) // Remove the phone from the map
				}
				return
			default:
				fmt.Println("Unhandled QR event:", evt.Event)
			}
		}
	}()
}

// CheckActivation checks if WhatsApp service is active for the given phone number
func (s *WaServiceImpl) Deactivate(ctx context.Context, phone string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if status exists and deactivate the service
	if status, exists := s.UserStatusMap[phone]; exists {
		status.IsActive = false
		status.WaClient.Disconnect()
		err := status.Container.Close()
		helper.PanicIfError("Error closing SQL container", err)
		// Here, cancel any active context or Goroutines related to this phone number
		return true
	}
	return false
}

func (s *WaServiceImpl) CheckActivation(ctx context.Context, phone string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if status exists and return its IsActive status
	if status, exists := s.UserStatusMap[phone]; exists {
		return status.IsActive
	}
	return false
}

// CheckAuthentication checks if WhatsApp service is authenticated for the given phone number
func (s *WaServiceImpl) CheckAuthentication(ctx context.Context, phone string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if status exists and return its IsAuthenticated status
	if status, exists := s.UserStatusMap[phone]; exists {
		return status.IsAuthenticated
	}
	return false
}
