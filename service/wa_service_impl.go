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
	mu            sync.Mutex               // Mutex to protect concurrent access to map
}

// Create new WhatsApp and AI service
func NewWaService() WaService {
	return &WaServiceImpl{
		UserStatusMap: make(map[string]*UserWaStatus),
	}
}

// Function to activate WhatsApp and AI for a user
func (s *WaServiceImpl) Activate(ctx context.Context, phone string) *UserWaStatus {
	// If the service is already active, no need to activate again
	s.mu.Lock()
	defer s.mu.Unlock()
	if status, exists := s.UserStatusMap[phone]; exists && status.IsActive {
		return status
	}

	waClient, container := app.GetWaClient(phone)
	waClient.Connect()
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
		s.handleQRCodeAuthentication(ctx, phone, status)
	}

	return status
}

// handleQRCodeAuthentication manages the WhatsApp authentication process using QR code
func (s *WaServiceImpl) handleQRCodeAuthentication(ctx context.Context, phone string, status *UserWaStatus) {
	qrChan, _ := status.WaClient.GetQRChannel(ctx)
	fmt.Println("Hanlde Qr Code")
	go func() {
		for evt := range qrChan {
			switch evt.Event {
			case "code":
				err := qrcode.WriteFile(evt.Code, qrcode.Medium, 256, "public/qr"+phone+".png")
				if err != nil {
				} else {
					fmt.Printf("QR code generated for phone %s. Scan it using WhatsApp!\n", phone)
				}
			case "success":
				s.mu.Lock()
				status.IsAuthenticated = true
				status.WaClient.Connect()
				s.mu.Unlock()
				fmt.Println("WhatsApp authentication successful!")
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
