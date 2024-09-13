package service

import (
	"adiubaidah/adi-bot/app"
	"context"
	"fmt"
	"sync"

	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
)

// Struct for storing user service status
type UserWaStatus struct {
	WaClient        *whatsmeow.Client
	IsActive        bool
	IsAuthenticated bool
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
func (s *WaServiceImpl) Activate(ctx context.Context, phone string) bool {
	// If the service is already active, no need to activate again
	s.mu.Lock()
	defer s.mu.Unlock()
	if status, exists := s.UserStatusMap[phone]; exists && status.IsActive {
		return status.IsAuthenticated
	}

	// Initialize a new status if it doesn't exist
	s.UserStatusMap[phone] = &UserWaStatus{
		WaClient:        app.GetWaClient(phone),
		IsActive:        false,
		IsAuthenticated: false,
	}

	// Set the current status variable for convenience
	status := s.UserStatusMap[phone]
	status.IsAuthenticated = status.WaClient.Store.ID != nil
	fmt.Println("Wa Status", status.IsAuthenticated)

	// If not authenticated, handle the authentication process
	if !status.IsAuthenticated {
		s.handleQRCodeAuthentication(ctx, phone, status)
	}

	// Mark service as active once authenticated
	if status.IsAuthenticated {
		status.IsActive = true
	}

	return status.IsAuthenticated
}

// handleQRCodeAuthentication manages the WhatsApp authentication process using QR code
func (s *WaServiceImpl) handleQRCodeAuthentication(ctx context.Context, phone string, status *UserWaStatus) {
	qrChan, _ := status.WaClient.GetQRChannel(ctx)

	go func() {
		for evt := range qrChan {
			switch evt.Event {
			case "code":
				err := qrcode.WriteFile(evt.Code, qrcode.Medium, 256, "public/qr"+phone+".png")
				if err != nil {
					fmt.Println("Failed to generate QR code:", err)
				} else {
					fmt.Printf("QR code generated for phone %s. Scan it using WhatsApp!\n", phone)
				}
			case "success":
				s.mu.Lock()
				status.IsAuthenticated = true
				s.mu.Unlock()
				fmt.Println("WhatsApp authentication successful!")
			default:
				fmt.Println("Unhandled QR event:", evt.Event)
			}
		}
	}()
}

func (s *WaServiceImpl) Deactivate(ctx context.Context, phone string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if status exists and deactivate the service
	if status, exists := s.UserStatusMap[phone]; exists {
		status.IsActive = false
		return true
	}
	return false
}

// CheckActivation checks if WhatsApp service is active for the given phone number
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
