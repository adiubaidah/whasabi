package service

import (
	"context"
)

type WaService interface {
	Activate(ctx context.Context, phone string) bool
	CheckActivation(ctx context.Context, phone string) bool
	CheckAuthentication(ctx context.Context, phone string) bool
	Deactivate(ctx context.Context, phone string) bool
}
