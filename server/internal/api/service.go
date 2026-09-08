package api

import "log/slog"

func NewAuthHandler(service AuthService, logger *slog.Logger, secureCookie bool) *AuthHandler {
	return &AuthHandler{service: service, logger: logger, secureCookie: secureCookie}
}
