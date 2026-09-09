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
	router.Put("/parents/:parentId/students/:studentId", h.GrantParentStudentAccess)
	router.Delete("/parents/:parentId/students/:studentId", h.RevokeParentStudentAccess)
	router.Get("/parents/:parentId/students", h.ListManagedParentStudents)
}
func (h *AcademicYearHandler) Register(router fiber.Router) {
	router.Get("/academic-years", h.ListAcademicYears)
	router.Post("/academic-years", h.CreateAcademicYear)
}
func (h *SubjectHandler) Register(router fiber.Router) {
	router.Get("/subjects", h.ListSubjects)
	router.Post("/subjects", h.CreateSubject)
}
func (h *ClassHandler) Register(router fiber.Router) {
	router.Get("/classes", h.ListClasses)
	router.Post("/classes", h.CreateClass)
	router.Get("/classes/:classId", h.GetClass)
	router.Get("/classes/:classId/students", h.ListClassStudents)
	router.Post("/classes/:classId/students", h.AddStudentToClass)
	router.Patch("/classes/:classId/students/:studentId", h.UpdateClassStudent)
	router.Get("/classes/:classId/subjects", h.ListClassSubjects)
	router.Post("/classes/:classId/subjects", h.AssignSubjectToClass)
	router.Patch("/class-subjects/:classSubjectId", h.UpdateClassSubject)
}
