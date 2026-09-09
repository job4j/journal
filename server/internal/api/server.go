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
	GrantParentStudentAccess(context.Context, string, uuid.UUID, uuid.UUID) (domain.ParentStudentAccess, error)
	RevokeParentStudentAccess(context.Context, string, uuid.UUID, uuid.UUID) (domain.ParentStudentAccess, error)
	ListManagedParentStudents(context.Context, string, uuid.UUID) ([]entity.User, error)
	ListParentStudents(context.Context, string) ([]entity.User, error)
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
	UpdateClassStudent(context.Context, string, uuid.UUID, uuid.UUID, *time.Time) (domain.ClassStudentView, error)
	ListClassSubjects(context.Context, string, uuid.UUID) ([]domain.ClassSubjectView, error)
	AssignSubjectToClass(context.Context, string, uuid.UUID, uuid.UUID, uuid.UUID) (domain.ClassSubjectView, error)
	UpdateClassSubject(context.Context, string, uuid.UUID, uuid.UUID) (domain.ClassSubjectView, error)
	ListTeacherClasses(context.Context, string) ([]domain.ClassView, error)
	ListTeacherClassSubjects(context.Context, string, uuid.UUID) ([]domain.ClassSubjectView, error)
	ListClassSubjectLessons(context.Context, string, uuid.UUID, *time.Time, *time.Time) ([]domain.LessonView, error)
	CreateLesson(context.Context, string, uuid.UUID, domain.CreateLessonInput) (domain.LessonView, error)
	CreateGradeItem(context.Context, string, uuid.UUID, string, string, string, *float64) (entity.GradeItem, error)
	PutStudentScore(context.Context, string, uuid.UUID, uuid.UUID, *float64, *string, *string) (entity.Score, error)
	PutStudentAbsence(context.Context, string, uuid.UUID, uuid.UUID) (entity.Absence, error)
	DeleteStudentAbsence(context.Context, string, uuid.UUID, uuid.UUID) error
	GetParentStudentJournal(context.Context, string, uuid.UUID, uuid.UUID) (domain.ParentJournal, error)
}
