package api

import "github.com/gofiber/fiber/v2"

func (h *AuthHandler) Register(router fiber.Router) { router.Post("/auth/login", h.Login) }
func (h *RoleHandler) Register(router fiber.Router) {
	router.Get("/roles", h.ListRoles)
	router.Post("/roles", h.CreateRole)
	router.Put("/roles/:roleId", h.UpdateRole)
	router.Delete("/roles/:roleId", h.DeleteRole)
	router.Get("/permissions", h.ListPermissions)
}
func (h *UserHandler) Register(router fiber.Router) {
	router.Get("/users", h.ListUsers)
	router.Post("/users", h.CreateUser)
	router.Get("/users/:userId", h.GetUser)
	router.Put("/users/:userId", h.UpdateUser)
	router.Delete("/users/:userId", h.DeleteUser)
}
