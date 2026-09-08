package api

import (
	"context"
	"journal/server/internal/domain"
	"log/slog"
)

type AuthHandler struct {
	service      AuthService
	logger       *slog.Logger
	secureCookie bool
}
type AuthService interface {
	Login(context.Context, string, string) (domain.LoginResponse, error)
}
