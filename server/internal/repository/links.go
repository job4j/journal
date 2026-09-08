package repository

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func (r *Repository) CreateUserRole(ctx context.Context, tx Transaction, value entity.UserRole) (entity.UserRole, error) {
	return queryOne[entity.UserRole](ctx, tx, "create user role", `INSERT INTO user_roles (user_id, role_id) VALUES ($1,$2) RETURNING user_id, role_id, created_at`, value.UserID, value.RoleID)
}
func (r *Repository) GetUserRole(ctx context.Context, tx Transaction, userID, roleID uuid.UUID) (entity.UserRole, error) {
	return queryOne[entity.UserRole](ctx, tx, "get user role", `SELECT user_id, role_id, created_at FROM user_roles WHERE user_id=$1 AND role_id=$2`, userID, roleID)
}
func (r *Repository) ListUserRoles(ctx context.Context, tx Transaction) ([]entity.UserRole, error) {
	return queryMany[entity.UserRole](ctx, tx, "list user roles", `SELECT user_id, role_id, created_at FROM user_roles ORDER BY created_at, user_id, role_id`)
}
func (r *Repository) DeleteUserRole(ctx context.Context, tx Transaction, userID, roleID uuid.UUID) error {
	return deleteRows(ctx, tx, "delete user role", `DELETE FROM user_roles WHERE user_id=$1 AND role_id=$2`, userID, roleID)
}

func (r *Repository) CreateRolePermission(ctx context.Context, tx Transaction, value entity.RolePermission) (entity.RolePermission, error) {
	return queryOne[entity.RolePermission](ctx, tx, "create role permission", `INSERT INTO role_permissions (role_id,permission_id) VALUES ($1,$2) RETURNING role_id,permission_id,created_at`, value.RoleID, value.PermissionID)
}
func (r *Repository) GetRolePermission(ctx context.Context, tx Transaction, roleID, permissionID uuid.UUID) (entity.RolePermission, error) {
	return queryOne[entity.RolePermission](ctx, tx, "get role permission", `SELECT role_id,permission_id,created_at FROM role_permissions WHERE role_id=$1 AND permission_id=$2`, roleID, permissionID)
}
func (r *Repository) ListRolePermissions(ctx context.Context, tx Transaction) ([]entity.RolePermission, error) {
	return queryMany[entity.RolePermission](ctx, tx, "list role permissions", `SELECT role_id,permission_id,created_at FROM role_permissions ORDER BY created_at,role_id,permission_id`)
}
func (r *Repository) DeleteRolePermission(ctx context.Context, tx Transaction, roleID, permissionID uuid.UUID) error {
	return deleteRows(ctx, tx, "delete role permission", `DELETE FROM role_permissions WHERE role_id=$1 AND permission_id=$2`, roleID, permissionID)
}

func (r *Repository) CreateUserPermission(ctx context.Context, tx Transaction, value entity.UserPermission) (entity.UserPermission, error) {
	return queryOne[entity.UserPermission](ctx, tx, "create user permission", `INSERT INTO user_permissions (user_id,permission_id) VALUES ($1,$2) RETURNING user_id,permission_id,created_at`, value.UserID, value.PermissionID)
}
func (r *Repository) GetUserPermission(ctx context.Context, tx Transaction, userID, permissionID uuid.UUID) (entity.UserPermission, error) {
	return queryOne[entity.UserPermission](ctx, tx, "get user permission", `SELECT user_id,permission_id,created_at FROM user_permissions WHERE user_id=$1 AND permission_id=$2`, userID, permissionID)
}
func (r *Repository) ListUserPermissions(ctx context.Context, tx Transaction) ([]entity.UserPermission, error) {
	return queryMany[entity.UserPermission](ctx, tx, "list user permissions", `SELECT user_id,permission_id,created_at FROM user_permissions ORDER BY created_at,user_id,permission_id`)
}
func (r *Repository) DeleteUserPermission(ctx context.Context, tx Transaction, userID, permissionID uuid.UUID) error {
	return deleteRows(ctx, tx, "delete user permission", `DELETE FROM user_permissions WHERE user_id=$1 AND permission_id=$2`, userID, permissionID)
}

func (r *Repository) CreateClassStudent(ctx context.Context, tx Transaction, value entity.ClassStudent) (entity.ClassStudent, error) {
	return queryOne[entity.ClassStudent](ctx, tx, "create class student", `INSERT INTO class_students (class_id,user_id,enrolled_on,left_on) VALUES ($1,$2,$3,$4) RETURNING class_id,user_id,enrolled_on,left_on,created_at,updated_at`, value.ClassID, value.UserID, value.EnrolledOn, value.LeftOn)
}
func (r *Repository) GetClassStudent(ctx context.Context, tx Transaction, classID, userID uuid.UUID) (entity.ClassStudent, error) {
	return queryOne[entity.ClassStudent](ctx, tx, "get class student", `SELECT class_id,user_id,enrolled_on,left_on,created_at,updated_at FROM class_students WHERE class_id=$1 AND user_id=$2`, classID, userID)
}
func (r *Repository) ListClassStudents(ctx context.Context, tx Transaction) ([]entity.ClassStudent, error) {
	return queryMany[entity.ClassStudent](ctx, tx, "list class students", `SELECT class_id,user_id,enrolled_on,left_on,created_at,updated_at FROM class_students ORDER BY created_at,class_id,user_id`)
}
func (r *Repository) UpdateClassStudent(ctx context.Context, tx Transaction, value entity.ClassStudent) (entity.ClassStudent, error) {
	return queryOne[entity.ClassStudent](ctx, tx, "update class student", `UPDATE class_students SET enrolled_on=$1,left_on=$2,updated_at=now() WHERE class_id=$3 AND user_id=$4 RETURNING class_id,user_id,enrolled_on,left_on,created_at,updated_at`, value.EnrolledOn, value.LeftOn, value.ClassID, value.UserID)
}
func (r *Repository) DeleteClassStudent(ctx context.Context, tx Transaction, classID, userID uuid.UUID) error {
	return deleteRows(ctx, tx, "delete class student", `DELETE FROM class_students WHERE class_id=$1 AND user_id=$2`, classID, userID)
}

func (r *Repository) ListRolePermissionsByRoleID(ctx context.Context, tx Transaction, roleID uuid.UUID) ([]entity.RolePermission, error) {
	return queryMany[entity.RolePermission](ctx, tx, "list role permissions by role", `SELECT role_id,permission_id,created_at FROM role_permissions WHERE role_id=$1 ORDER BY created_at,permission_id`, roleID)
}

func (r *Repository) ListUserRolesByUserID(ctx context.Context, tx Transaction, userID uuid.UUID) ([]entity.UserRole, error) {
	return queryMany[entity.UserRole](ctx, tx, "list user roles by user", `SELECT user_id,role_id,created_at FROM user_roles WHERE user_id=$1 ORDER BY created_at,role_id`, userID)
}
