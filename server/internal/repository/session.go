package repository

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func (r *Repository) CreateSession(ctx context.Context, tx Transaction, value entity.Session) (entity.Session, error) {
	return queryOne[entity.Session](ctx, tx, "create session", `INSERT INTO sessions (user_id, token_hash, expires_at, revoked_at) VALUES ($1, $2, $3, $4) RETURNING id, user_id, token_hash, expires_at, created_at, last_used_at, revoked_at`, value.UserID, value.TokenHash, value.ExpiresAt, value.RevokedAt)
}
func (r *Repository) GetSession(ctx context.Context, tx Transaction, id uuid.UUID) (entity.Session, error) {
	return queryOne[entity.Session](ctx, tx, "get session", `SELECT id, user_id, token_hash, expires_at, created_at, last_used_at, revoked_at FROM sessions WHERE id = $1`, id)
}
func (r *Repository) ListSessions(ctx context.Context, tx Transaction) ([]entity.Session, error) {
	return queryMany[entity.Session](ctx, tx, "list sessions", `SELECT id, user_id, token_hash, expires_at, created_at, last_used_at, revoked_at FROM sessions ORDER BY created_at, id`)
}
func (r *Repository) UpdateSession(ctx context.Context, tx Transaction, value entity.Session) (entity.Session, error) {
	return queryOne[entity.Session](ctx, tx, "update session", `UPDATE sessions SET user_id = $1, token_hash = $2, expires_at = $3, last_used_at = $4, revoked_at = $5 WHERE id = $6 RETURNING id, user_id, token_hash, expires_at, created_at, last_used_at, revoked_at`, value.UserID, value.TokenHash, value.ExpiresAt, value.LastUsedAt, value.RevokedAt, value.ID)
}
func (r *Repository) DeleteSession(ctx context.Context, tx Transaction, id uuid.UUID) error {
	return deleteRows(ctx, tx, "delete session", `DELETE FROM sessions WHERE id = $1`, id)
}
