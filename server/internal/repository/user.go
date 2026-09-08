package repository

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

const userColumns = `u.id,u.email,u.password_hash,u.first_name,u.last_name,u.status,
 COALESCE((SELECT array_agg(r.code ORDER BY r.code) FROM user_roles ur JOIN roles r ON r.id=ur.role_id WHERE ur.user_id=u.id),'{}') AS roles,
 u.created_at,u.updated_at`

func (r *Repository) CreateUser(ctx context.Context, tx Transaction, value entity.User) (entity.User, error) {
	return queryOne[entity.User](ctx, tx, "create user", `INSERT INTO users(email,password_hash,first_name,last_name,status) VALUES($1,$2,$3,$4,$5) RETURNING id,email,password_hash,first_name,last_name,status,'{}'::text[] AS roles,created_at,updated_at`, value.Email, value.PasswordHash, value.FirstName, value.LastName, value.Status)
}
func (r *Repository) GetUser(ctx context.Context, tx Transaction, id uuid.UUID) (entity.User, error) {
	return queryOne[entity.User](ctx, tx, "get user", `SELECT `+userColumns+` FROM users u WHERE u.id=$1`, id)
}
func (r *Repository) ListUsers(ctx context.Context, tx Transaction) ([]entity.User, error) {
	return queryMany[entity.User](ctx, tx, "list users", `SELECT `+userColumns+` FROM users u ORDER BY u.created_at,u.id`)
}
func (r *Repository) UpdateUser(ctx context.Context, tx Transaction, value entity.User) (entity.User, error) {
	return queryOne[entity.User](ctx, tx, "update user", `UPDATE users SET email=$1,password_hash=$2,first_name=$3,last_name=$4,status=$5,updated_at=now() WHERE id=$6 RETURNING id,email,password_hash,first_name,last_name,status,'{}'::text[] AS roles,created_at,updated_at`, value.Email, value.PasswordHash, value.FirstName, value.LastName, value.Status, value.ID)
}
func (r *Repository) DeleteUser(ctx context.Context, tx Transaction, id uuid.UUID) error {
	return deleteRows(ctx, tx, "delete user", `DELETE FROM users WHERE id=$1`, id)
}
