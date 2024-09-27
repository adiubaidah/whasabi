package service

import (
	"context"
)

type WaService interface {
	Activate(phone string) *UserWaStatus
	CheckActivation(ctx context.Context, phone string) bool
	CheckAuthentication(ctx context.Context, phone string) bool
	Deactivate(ctx context.Context, phone string) bool
}
