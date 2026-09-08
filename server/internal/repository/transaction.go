package repository

import "context"

type Transaction interface {
	Commit(context.Context) error
	Rollback(context.Context) error
}
type TransactionManager interface {
	Begin(context.Context) (Transaction, error)
}
