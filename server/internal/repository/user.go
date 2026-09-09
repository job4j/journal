package repository

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

const userColumns = `u.id,u.login,COALESCE(u.email,''),COALESCE(u.phone,''),u.password_hash,u.first_name,u.last_name,u.status,
 COALESCE((SELECT array_agg(r.code ORDER BY r.code) FROM user_roles ur JOIN roles r ON r.id=ur.role_id WHERE ur.user_id=u.id),'{}') AS roles,
 u.created_at,u.updated_at`

func (r *Repository) CreateUser(ctx context.Context, tx Transaction, value entity.User) (entity.User, error) {
	return queryOne[entity.User](ctx, tx, "create user", `INSERT INTO users(login,email,phone,password_hash,first_name,last_name,status) VALUES($1,NULLIF($2,''),NULLIF($3,''),$4,$5,$6,$7) RETURNING id,login,COALESCE(email,''),COALESCE(phone,''),password_hash,first_name,last_name,status,'{}'::text[] AS roles,created_at,updated_at`, value.Login, value.Email, value.Phone, value.PasswordHash, value.FirstName, value.LastName, value.Status)
}
func (r *Repository) GetUser(ctx context.Context, tx Transaction, id uuid.UUID) (entity.User, error) {
	return queryOne[entity.User](ctx, tx, "get user", `SELECT `+userColumns+` FROM users u WHERE u.id=$1`, id)
}
func (r *Repository) FindActiveUserBySessionHash(ctx context.Context, tx Transaction, hash string) (entity.User, error) {
	return queryOne[entity.User](ctx, tx, "find active user by session", `SELECT `+userColumns+` FROM users u JOIN sessions s ON s.user_id=u.id WHERE s.token_hash=$1 AND s.revoked_at IS NULL AND s.expires_at>now() AND u.status='active'`, hash)
}
func (r *Repository) ListUsers(ctx context.Context, tx Transaction) ([]entity.User, error) {
	return queryMany[entity.User](ctx, tx, "list users", `SELECT `+userColumns+` FROM users u ORDER BY u.created_at,u.id`)
}
func (r *Repository) UpdateUser(ctx context.Context, tx Transaction, value entity.User) (entity.User, error) {
	return queryOne[entity.User](ctx, tx, "update user", `UPDATE users SET login=$1,email=NULLIF($2,''),phone=NULLIF($3,''),password_hash=$4,first_name=$5,last_name=$6,status=$7,updated_at=now() WHERE id=$8 RETURNING id,login,COALESCE(email,''),COALESCE(phone,''),password_hash,first_name,last_name,status,'{}'::text[] AS roles,created_at,updated_at`, value.Login, value.Email, value.Phone, value.PasswordHash, value.FirstName, value.LastName, value.Status, value.ID)
}
func (r *Repository) DeleteUser(ctx context.Context, tx Transaction, id uuid.UUID) error {
	return deleteRows(ctx, tx, "delete user", `DELETE FROM users WHERE id=$1`, id)
}
