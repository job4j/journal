package api

import "github.com/gofiber/fiber/v2"

func (h *AuthHandler) Register(router fiber.Router) { router.Post("/auth/login", h.Login) }
