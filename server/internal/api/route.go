package api

import "github.com/gofiber/fiber/v2"

func (h *AuthHandler) Register(router fiber.Router) {
	router.Post("/auth/login", h.Login)
	router.Post("/auth/logout", h.Logout)
	router.Get("/me", h.CurrentUser)
}
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
func (h *AcademicYearHandler) Register(router fiber.Router) {
	router.Get("/academic-years", h.ListAcademicYears)
	router.Post("/academic-years", h.CreateAcademicYear)
}
func (h *SubjectHandler) Register(router fiber.Router) {
	router.Get("/subjects", h.ListSubjects)
	router.Post("/subjects", h.CreateSubject)
}
