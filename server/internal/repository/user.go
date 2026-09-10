package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

const userColumns = `u.id,u.login,u.email,u.phone,u.password_hash,u.first_name,u.last_name,u.status,
 ARRAY(SELECT r.code FROM user_roles ur JOIN roles r ON r.id=ur.role_id WHERE ur.user_id=u.id ORDER BY r.code) AS roles,
 u.created_at,u.updated_at`

type userRow struct {
	ID           uuid.UUID `db:"id"`
	Login        string    `db:"login"`
	Email        *string   `db:"email"`
	Phone        *string   `db:"phone"`
	PasswordHash string    `db:"password_hash"`
	FirstName    string    `db:"first_name"`
	LastName     string    `db:"last_name"`
	Status       string    `db:"status"`
	Roles        []string  `db:"roles"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (row userRow) entity() entity.User {
	user := entity.User{ID: row.ID, Login: row.Login, PasswordHash: row.PasswordHash, FirstName: row.FirstName, LastName: row.LastName, Status: row.Status, Roles: row.Roles, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt}
	if row.Email != nil {
		user.Email = *row.Email
	}
	if row.Phone != nil {
		user.Phone = *row.Phone
	}
	return user
}

func (r *Repository) CreateUser(ctx context.Context, tx Transaction, value entity.User) (entity.User, error) {
	row, err := queryOne[userRow](ctx, tx, "create user", `INSERT INTO users(login,email,phone,password_hash,first_name,last_name,status) VALUES($1,NULLIF($2,''),NULLIF($3,''),$4,$5,$6,$7) RETURNING id,login,email,phone,password_hash,first_name,last_name,status,'{}'::text[] AS roles,created_at,updated_at`, value.Login, value.Email, value.Phone, value.PasswordHash, value.FirstName, value.LastName, value.Status)
	return row.entity(), err
}
func (r *Repository) GetUser(ctx context.Context, tx Transaction, id uuid.UUID) (entity.User, error) {
	row, err := queryOne[userRow](ctx, tx, "get user", `SELECT `+userColumns+` FROM users u WHERE u.id=$1`, id)
	return row.entity(), err
}
func (r *Repository) FindActiveUserBySessionHash(ctx context.Context, tx Transaction, hash string) (entity.User, error) {
	row, err := queryOne[userRow](ctx, tx, "find active user by session", `SELECT `+userColumns+` FROM users u JOIN sessions s ON s.user_id=u.id WHERE s.token_hash=$1 AND s.revoked_at IS NULL AND s.expires_at>now() AND u.status='active'`, hash)
	return row.entity(), err
}
func (r *Repository) ListUsers(ctx context.Context, tx Transaction) ([]entity.User, error) {
	rows, err := queryMany[userRow](ctx, tx, "list users", `SELECT `+userColumns+` FROM users u ORDER BY u.created_at,u.id`)
	if err != nil {
		return nil, err
	}
	users := make([]entity.User, len(rows))
	for i := range rows {
		users[i] = rows[i].entity()
	}
	return users, nil
}
func (r *Repository) UpdateUser(ctx context.Context, tx Transaction, value entity.User) (entity.User, error) {
	row, err := queryOne[userRow](ctx, tx, "update user", `UPDATE users SET login=$1,email=NULLIF($2,''),phone=NULLIF($3,''),password_hash=$4,first_name=$5,last_name=$6,status=$7,updated_at=now() WHERE id=$8 RETURNING id,login,email,phone,password_hash,first_name,last_name,status,'{}'::text[] AS roles,created_at,updated_at`, value.Login, value.Email, value.Phone, value.PasswordHash, value.FirstName, value.LastName, value.Status, value.ID)
	return row.entity(), err
}
func (r *Repository) DeleteUser(ctx context.Context, tx Transaction, id uuid.UUID) error {
	return deleteRows(ctx, tx, "delete user", `DELETE FROM users WHERE id=$1`, id)
}
