package repository

import (
	"context"

	"journal/server/internal/repository/entity"
)

type authRepository struct{}

func NewAuthRepository() *authRepository { return &authRepository{} }

type AuthRepository interface {
	FindUserByEmail(context.Context, Transaction, string) (entity.User, error)
	InsertSession(context.Context, Transaction, entity.Session) error
}
