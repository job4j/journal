package api

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/domain"
	"journal/server/internal/repository/entity"
	"log/slog"
	"time"
)

type AuthHandler struct {
	service      AuthService
	logger       *slog.Logger
	secureCookie bool
}
type AuthService interface {
	Login(context.Context, string, string) (domain.LoginResponse, error)
	CurrentUser(context.Context, string) (entity.User, error)
	Logout(context.Context, string) error
}
type RoleHandler struct {
	service RoleService
	logger  *slog.Logger
}
type RoleService interface {
	CreateRole(context.Context, string, string, string, []string) (domain.CreateRoleResponse, error)
	ListRoles(context.Context, string) (domain.ListRolesResponse, error)
	ListPermissions(context.Context, string) (domain.ListPermissionsResponse, error)
	UpdateRole(context.Context, string, uuid.UUID, string, string, []string) (domain.RoleView, error)
	DeleteRole(context.Context, string, uuid.UUID) error
}
type UserHandler struct {
	service UserService
	logger  *slog.Logger
}
type UserService interface {
	CreateUser(context.Context, string, domain.UserInput) (entity.User, error)
	ListUsers(context.Context, string, string, int, int) (domain.ListUsersResult, error)
	GetUser(context.Context, string, uuid.UUID) (entity.User, error)
	UpdateUser(context.Context, string, uuid.UUID, domain.UserInput) (entity.User, error)
	DeleteUser(context.Context, string, uuid.UUID) error
}
type AcademicYearHandler struct {
	service AcademicYearService
	logger  *slog.Logger
}
type AcademicYearService interface {
	ListAcademicYears(context.Context, string) ([]domain.AcademicYearView, error)
	CreateAcademicYear(context.Context, string, domain.CreateAcademicYearRequest) (domain.AcademicYearView, error)
}
type SubjectHandler struct {
	service SubjectService
	logger  *slog.Logger
}
type SubjectService interface {
	ListSubjects(context.Context, string) ([]entity.Subject, error)
	CreateSubject(context.Context, string, string, string) (entity.Subject, error)
}
type ClassHandler struct {
	service ClassService
	logger  *slog.Logger
}
type ClassService interface {
	ListClasses(context.Context, string, uuid.UUID) ([]domain.ClassView, error)
	GetClass(context.Context, string, uuid.UUID) (domain.ClassView, error)
	CreateClass(context.Context, string, uuid.UUID, string, int16) (domain.ClassView, error)
	ListClassStudents(context.Context, string, uuid.UUID) ([]domain.ClassStudentView, error)
	AddStudentToClass(context.Context, string, uuid.UUID, uuid.UUID, time.Time) (domain.ClassStudentView, error)
}
