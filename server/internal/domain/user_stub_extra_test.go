package domain

import (
	"context"
	"journal/server/internal/repository"
)

func (r *userRepoStub) CheckSessionPermissionForValue(context.Context, repository.Transaction, string, string, string) (bool, bool, error) {
	return r.authenticated, r.allowed, nil
}
