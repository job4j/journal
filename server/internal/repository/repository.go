package repository

import "context"

type AuthRepository interface {
	FindUserByEmail(context.Context, Transaction, string) (User, error)
	InsertSession(context.Context, Transaction, Session) error
}
