package api

import "log/slog"

func NewAuthHandler(service AuthService, logger *slog.Logger, secureCookie bool) *AuthHandler {
	return &AuthHandler{service: service, logger: logger, secureCookie: secureCookie}
}
func NewRoleHandler(service RoleService, logger *slog.Logger) *RoleHandler {
	return &RoleHandler{service: service, logger: logger}
}
func NewUserHandler(service UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{service: service, logger: logger}
}
func NewAcademicYearHandler(service AcademicYearService, logger *slog.Logger) *AcademicYearHandler {
	return &AcademicYearHandler{service: service, logger: logger}
}
func NewSubjectHandler(service SubjectService, logger *slog.Logger) *SubjectHandler {
	return &SubjectHandler{service: service, logger: logger}
}
