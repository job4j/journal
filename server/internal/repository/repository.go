package repository

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
	"time"
)

type authRepository struct{}

func NewAuthRepository() *authRepository { return &authRepository{} }

type AuthRepository interface {
	FindUserByEmail(context.Context, Transaction, string) (entity.User, error)
	FindUserBySessionTokenHash(context.Context, Transaction, string, time.Time) (entity.User, error)
	InsertSession(context.Context, Transaction, entity.Session) error
	RevokeSession(context.Context, Transaction, string, time.Time) error
}

type UserRepository interface {
	CreateUser(context.Context, Transaction, entity.User) (entity.User, error)
	GetUser(context.Context, Transaction, uuid.UUID) (entity.User, error)
	ListUsers(context.Context, Transaction) ([]entity.User, error)
	UpdateUser(context.Context, Transaction, entity.User) (entity.User, error)
	DeleteUser(context.Context, Transaction, uuid.UUID) error
}
type AccessRepository interface {
	CreateRole(context.Context, Transaction, entity.Role) (entity.Role, error)
	GetRole(context.Context, Transaction, uuid.UUID) (entity.Role, error)
	ListRoles(context.Context, Transaction) ([]entity.Role, error)
	UpdateRole(context.Context, Transaction, entity.Role) (entity.Role, error)
	DeleteRole(context.Context, Transaction, uuid.UUID) error
	CreatePermission(context.Context, Transaction, entity.Permission) (entity.Permission, error)
	GetPermission(context.Context, Transaction, uuid.UUID) (entity.Permission, error)
	ListPermissions(context.Context, Transaction) ([]entity.Permission, error)
	UpdatePermission(context.Context, Transaction, entity.Permission) (entity.Permission, error)
	DeletePermission(context.Context, Transaction, uuid.UUID) error
	CreateUserRole(context.Context, Transaction, entity.UserRole) (entity.UserRole, error)
	GetUserRole(context.Context, Transaction, uuid.UUID, uuid.UUID) (entity.UserRole, error)
	ListUserRoles(context.Context, Transaction) ([]entity.UserRole, error)
	DeleteUserRole(context.Context, Transaction, uuid.UUID, uuid.UUID) error
	CreateRolePermission(context.Context, Transaction, entity.RolePermission) (entity.RolePermission, error)
	GetRolePermission(context.Context, Transaction, uuid.UUID, uuid.UUID) (entity.RolePermission, error)
	ListRolePermissions(context.Context, Transaction) ([]entity.RolePermission, error)
	DeleteRolePermission(context.Context, Transaction, uuid.UUID, uuid.UUID) error
	CreateUserPermission(context.Context, Transaction, entity.UserPermission) (entity.UserPermission, error)
	GetUserPermission(context.Context, Transaction, uuid.UUID, uuid.UUID) (entity.UserPermission, error)
	ListUserPermissions(context.Context, Transaction) ([]entity.UserPermission, error)
	DeleteUserPermission(context.Context, Transaction, uuid.UUID, uuid.UUID) error
}
type AcademicRepository interface {
	CreateAcademicYear(context.Context, Transaction, entity.AcademicYear) (entity.AcademicYear, error)
	GetAcademicYear(context.Context, Transaction, uuid.UUID) (entity.AcademicYear, error)
	ListAcademicYears(context.Context, Transaction) ([]entity.AcademicYear, error)
	UpdateAcademicYear(context.Context, Transaction, entity.AcademicYear) (entity.AcademicYear, error)
	DeleteAcademicYear(context.Context, Transaction, uuid.UUID) error
	CreateAcademicYearQuarter(context.Context, Transaction, entity.AcademicYearQuarter) (entity.AcademicYearQuarter, error)
	GetAcademicYearQuarter(context.Context, Transaction, uuid.UUID) (entity.AcademicYearQuarter, error)
	ListAcademicYearQuarters(context.Context, Transaction) ([]entity.AcademicYearQuarter, error)
	UpdateAcademicYearQuarter(context.Context, Transaction, entity.AcademicYearQuarter) (entity.AcademicYearQuarter, error)
	DeleteAcademicYearQuarter(context.Context, Transaction, uuid.UUID) error
}
type AcademicYearRepository interface {
	CheckSessionPermission(context.Context, Transaction, string, string) (bool, bool, error)
	CreateAcademicYear(context.Context, Transaction, entity.AcademicYear) (entity.AcademicYear, error)
	CreateAcademicYearQuarter(context.Context, Transaction, entity.AcademicYearQuarter) (entity.AcademicYearQuarter, error)
	ListAcademicYears(context.Context, Transaction) ([]entity.AcademicYear, error)
	ListAcademicYearQuarters(context.Context, Transaction) ([]entity.AcademicYearQuarter, error)
}
type SubjectRepository interface {
	CheckSessionPermission(context.Context, Transaction, string, string) (bool, bool, error)
	CreateSubject(context.Context, Transaction, entity.Subject) (entity.Subject, error)
	ListSubjects(context.Context, Transaction) ([]entity.Subject, error)
}
type ClassRepository interface {
	CreateClass(context.Context, Transaction, entity.Class) (entity.Class, error)
	GetClass(context.Context, Transaction, uuid.UUID) (entity.Class, error)
	ListClasses(context.Context, Transaction) ([]entity.Class, error)
	UpdateClass(context.Context, Transaction, entity.Class) (entity.Class, error)
	DeleteClass(context.Context, Transaction, uuid.UUID) error
	CreateSubject(context.Context, Transaction, entity.Subject) (entity.Subject, error)
	GetSubject(context.Context, Transaction, uuid.UUID) (entity.Subject, error)
	ListSubjects(context.Context, Transaction) ([]entity.Subject, error)
	UpdateSubject(context.Context, Transaction, entity.Subject) (entity.Subject, error)
	DeleteSubject(context.Context, Transaction, uuid.UUID) error
	CreateClassSubject(context.Context, Transaction, entity.ClassSubject) (entity.ClassSubject, error)
	GetClassSubject(context.Context, Transaction, uuid.UUID) (entity.ClassSubject, error)
	ListClassSubjects(context.Context, Transaction) ([]entity.ClassSubject, error)
	UpdateClassSubject(context.Context, Transaction, entity.ClassSubject) (entity.ClassSubject, error)
	DeleteClassSubject(context.Context, Transaction, uuid.UUID) error
	CreateClassStudent(context.Context, Transaction, entity.ClassStudent) (entity.ClassStudent, error)
	GetClassStudent(context.Context, Transaction, uuid.UUID, uuid.UUID) (entity.ClassStudent, error)
	ListClassStudents(context.Context, Transaction) ([]entity.ClassStudent, error)
	UpdateClassStudent(context.Context, Transaction, entity.ClassStudent) (entity.ClassStudent, error)
	DeleteClassStudent(context.Context, Transaction, uuid.UUID, uuid.UUID) error
}
type LessonRepository interface {
	CreateLesson(context.Context, Transaction, entity.Lesson) (entity.Lesson, error)
	GetLesson(context.Context, Transaction, uuid.UUID) (entity.Lesson, error)
	ListLessons(context.Context, Transaction) ([]entity.Lesson, error)
	UpdateLesson(context.Context, Transaction, entity.Lesson) (entity.Lesson, error)
	DeleteLesson(context.Context, Transaction, uuid.UUID) error
	CreateLessonMaterial(context.Context, Transaction, entity.LessonMaterial) (entity.LessonMaterial, error)
	GetLessonMaterial(context.Context, Transaction, uuid.UUID) (entity.LessonMaterial, error)
	ListLessonMaterials(context.Context, Transaction) ([]entity.LessonMaterial, error)
	UpdateLessonMaterial(context.Context, Transaction, entity.LessonMaterial) (entity.LessonMaterial, error)
	DeleteLessonMaterial(context.Context, Transaction, uuid.UUID) error
}
type GradingRepository interface {
	CreateGradeItem(context.Context, Transaction, entity.GradeItem) (entity.GradeItem, error)
	GetGradeItem(context.Context, Transaction, uuid.UUID) (entity.GradeItem, error)
	ListGradeItems(context.Context, Transaction) ([]entity.GradeItem, error)
	UpdateGradeItem(context.Context, Transaction, entity.GradeItem) (entity.GradeItem, error)
	DeleteGradeItem(context.Context, Transaction, uuid.UUID) error
	CreateScore(context.Context, Transaction, entity.Score) (entity.Score, error)
	GetScore(context.Context, Transaction, uuid.UUID) (entity.Score, error)
	ListScores(context.Context, Transaction) ([]entity.Score, error)
	UpdateScore(context.Context, Transaction, entity.Score) (entity.Score, error)
	DeleteScore(context.Context, Transaction, uuid.UUID) error
}
type SessionRepository interface {
	CreateSession(context.Context, Transaction, entity.Session) (entity.Session, error)
	GetSession(context.Context, Transaction, uuid.UUID) (entity.Session, error)
	ListSessions(context.Context, Transaction) ([]entity.Session, error)
	UpdateSession(context.Context, Transaction, entity.Session) (entity.Session, error)
	DeleteSession(context.Context, Transaction, uuid.UUID) error
}

type RoleRepository interface {
	CheckSessionPermission(context.Context, Transaction, string, string) (bool, bool, error)
	FindGlobalPermissionsByCodes(context.Context, Transaction, []string) ([]entity.Permission, error)
	ListGlobalPermissions(context.Context, Transaction) ([]entity.Permission, error)
	CreateRole(context.Context, Transaction, entity.Role) (entity.Role, error)
	GetRole(context.Context, Transaction, uuid.UUID) (entity.Role, error)
	ListRoles(context.Context, Transaction) ([]entity.Role, error)
	UpdateRole(context.Context, Transaction, entity.Role) (entity.Role, error)
	DeleteRole(context.Context, Transaction, uuid.UUID) error
	CreateRolePermission(context.Context, Transaction, entity.RolePermission) (entity.RolePermission, error)
	ListRolePermissions(context.Context, Transaction) ([]entity.RolePermission, error)
	ListRolePermissionsByRoleID(context.Context, Transaction, uuid.UUID) ([]entity.RolePermission, error)
	DeleteRolePermission(context.Context, Transaction, uuid.UUID, uuid.UUID) error
}

type UserManagementRepository interface {
	CheckSessionPermission(context.Context, Transaction, string, string) (bool, bool, error)
	CheckSessionPermissionForValue(context.Context, Transaction, string, string, string) (bool, bool, error)
	CreateUser(context.Context, Transaction, entity.User) (entity.User, error)
	GetUser(context.Context, Transaction, uuid.UUID) (entity.User, error)
	ListUsers(context.Context, Transaction) ([]entity.User, error)
	UpdateUser(context.Context, Transaction, entity.User) (entity.User, error)
	DeleteUser(context.Context, Transaction, uuid.UUID) error
	FindRolesByCodes(context.Context, Transaction, []string) ([]entity.Role, error)
	CreateUserRole(context.Context, Transaction, entity.UserRole) (entity.UserRole, error)
	ListUserRolesByUserID(context.Context, Transaction, uuid.UUID) ([]entity.UserRole, error)
	DeleteUserRole(context.Context, Transaction, uuid.UUID, uuid.UUID) error
}
