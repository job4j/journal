package repository

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func (r *Repository) CreatePermission(ctx context.Context, tx Transaction, value entity.Permission) (entity.Permission, error) {
	return queryOne[entity.Permission](ctx, tx, "create permission", `INSERT INTO permissions (code, value, description) VALUES ($1, $2, $3) RETURNING id, code, value, description, created_at`, value.Code, value.Value, value.Description)
}

func (r *Repository) EnsurePermission(
	ctx context.Context,
	tx Transaction,
	value entity.Permission,
) (entity.Permission, error) {
	if value.Value == nil {
		return queryOne[entity.Permission](ctx, tx, "ensure global permission", `
			INSERT INTO permissions (code, value, description)
			VALUES ($1, NULL, $2)
			ON CONFLICT (code) WHERE value IS NULL
			DO UPDATE SET description = EXCLUDED.description
			RETURNING id, code, value, description, created_at`,
			value.Code,
			value.Description,
		)
	}
	return queryOne[entity.Permission](ctx, tx, "ensure object permission", `
		INSERT INTO permissions (code, value, description)
		VALUES ($1, $2, $3)
		ON CONFLICT (code, value) WHERE value IS NOT NULL
		DO UPDATE SET description = EXCLUDED.description
		RETURNING id, code, value, description, created_at`,
		value.Code,
		value.Value,
		value.Description,
	)
}
func (r *Repository) GetPermission(ctx context.Context, tx Transaction, id uuid.UUID) (entity.Permission, error) {
	return queryOne[entity.Permission](ctx, tx, "get permission", `SELECT id, code, value, description, created_at FROM permissions WHERE id = $1`, id)
}
func (r *Repository) ListPermissions(ctx context.Context, tx Transaction) ([]entity.Permission, error) {
	return queryMany[entity.Permission](ctx, tx, "list permissions", `SELECT id, code, value, description, created_at FROM permissions ORDER BY created_at, id`)
}
func (r *Repository) UpdatePermission(ctx context.Context, tx Transaction, value entity.Permission) (entity.Permission, error) {
	return queryOne[entity.Permission](ctx, tx, "update permission", `UPDATE permissions SET code = $1, value = $2, description = $3 WHERE id = $4 RETURNING id, code, value, description, created_at`, value.Code, value.Value, value.Description, value.ID)
}
func (r *Repository) DeletePermission(ctx context.Context, tx Transaction, id uuid.UUID) error {
	return deleteRows(ctx, tx, "delete permission", `DELETE FROM permissions WHERE id = $1`, id)
}

func (r *Repository) FindGlobalPermissionsByCodes(ctx context.Context, tx Transaction, codes []string) ([]entity.Permission, error) {
	return queryMany[entity.Permission](ctx, tx, "find global permissions by codes", `SELECT id, code, value, description, created_at FROM permissions WHERE value IS NULL AND code = ANY($1) ORDER BY code`, codes)
}

func (r *Repository) CheckSessionPermission(ctx context.Context, transaction Transaction, tokenHash, permissionCode string) (bool, bool, error) {
	tx, err := pgxTransaction(transaction)
	if err != nil {
		return false, false, err
	}
	var authenticated, allowed bool
	err = tx.QueryRow(ctx, `
  SELECT
   EXISTS (
    SELECT 1 FROM sessions s
    JOIN users u ON u.id = s.user_id
    WHERE s.token_hash = $1 AND s.revoked_at IS NULL AND s.expires_at > now() AND u.status = 'active'
   ),
   EXISTS (
    SELECT 1 FROM sessions s
    JOIN users u ON u.id = s.user_id
    WHERE s.token_hash = $1 AND s.revoked_at IS NULL AND s.expires_at > now() AND u.status = 'active'
      AND (
       EXISTS (
        SELECT 1 FROM user_permissions up
        JOIN permissions p ON p.id = up.permission_id
        WHERE up.user_id = u.id AND p.code = $2 AND p.value IS NULL
       )
       OR EXISTS (
        SELECT 1 FROM user_roles ur
        JOIN role_permissions rp ON rp.role_id = ur.role_id
        JOIN permissions p ON p.id = rp.permission_id
        WHERE ur.user_id = u.id AND p.code = $2 AND p.value IS NULL
       )
      )
   )
 `, tokenHash, permissionCode).Scan(&authenticated, &allowed)
	if err != nil {
		return false, false, operationError("check session permission", err)
	}
	return authenticated, allowed, nil
}

func (r *Repository) ListGlobalPermissions(ctx context.Context, tx Transaction) ([]entity.Permission, error) {
	return queryMany[entity.Permission](ctx, tx, "list global permissions", `SELECT id,code,value,description,created_at FROM permissions WHERE value IS NULL ORDER BY code`)
}

func (r *Repository) CheckSessionPermissionForValue(ctx context.Context, transaction Transaction, tokenHash, permissionCode, value string) (bool, bool, error) {
	tx, err := pgxTransaction(transaction)
	if err != nil {
		return false, false, err
	}
	var authenticated, allowed bool
	err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM sessions s JOIN users u ON u.id=s.user_id WHERE s.token_hash=$1 AND s.revoked_at IS NULL AND s.expires_at>now() AND u.status='active'),EXISTS(SELECT 1 FROM sessions s JOIN users u ON u.id=s.user_id WHERE s.token_hash=$1 AND s.revoked_at IS NULL AND s.expires_at>now() AND u.status='active' AND(EXISTS(SELECT 1 FROM user_permissions up JOIN permissions p ON p.id=up.permission_id WHERE up.user_id=u.id AND p.code=$2 AND(p.value IS NULL OR p.value=$3))OR EXISTS(SELECT 1 FROM user_roles ur JOIN role_permissions rp ON rp.role_id=ur.role_id JOIN permissions p ON p.id=rp.permission_id WHERE ur.user_id=u.id AND p.code=$2 AND(p.value IS NULL OR p.value=$3))))`, tokenHash, permissionCode, value).Scan(&authenticated, &allowed)
	if err != nil {
		return false, false, operationError("check session permission value", err)
	}
	return authenticated, allowed, nil
}
