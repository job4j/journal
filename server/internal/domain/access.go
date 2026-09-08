package domain

import (
	"context"
	"fmt"
	"journal/server/internal/repository"
)

type AccessChecker interface {
	CheckSessionPermission(context.Context, repository.Transaction, string, string) (bool, bool, error)
}

type ObjectAccessChecker interface {
	CheckSessionPermissionForValue(context.Context, repository.Transaction, string, string, string) (bool, bool, error)
}

func Authorize(ctx context.Context, tx repository.Transaction, checker AccessChecker, tokenHash, permission string) error {
	authenticated, allowed, err := checker.CheckSessionPermission(ctx, tx, tokenHash, permission)
	return accessResult(authenticated, allowed, err)
}

func AuthorizeObject(ctx context.Context, tx repository.Transaction, checker ObjectAccessChecker, tokenHash, permission, value string) error {
	authenticated, allowed, err := checker.CheckSessionPermissionForValue(ctx, tx, tokenHash, permission, value)
	return accessResult(authenticated, allowed, err)
}

func accessResult(authenticated, allowed bool, err error) error {
	if err != nil {
		return fmt.Errorf("check access: %w", err)
	}
	if !authenticated {
		return ErrUnauthenticated
	}
	if !allowed {
		return ErrForbidden
	}
	return nil
}
