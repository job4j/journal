package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"testing"
)

type roleTransaction struct {
	committed  bool
	rolledBack bool
}

func (t *roleTransaction) Commit(context.Context) error   { t.committed = true; return nil }
func (t *roleTransaction) Rollback(context.Context) error { t.rolledBack = true; return nil }

type roleTransactionManager struct{ tx *roleTransaction }

func (m roleTransactionManager) Begin(context.Context) (repository.Transaction, error) {
	return m.tx, nil
}

type serviceRoleRepoStub struct {
	tokenHash string
	assignErr error
}

func (r *serviceRoleRepoStub) CheckSessionPermission(_ context.Context, _ repository.Transaction, tokenHash, _ string) (bool, bool, error) {
	r.tokenHash = tokenHash
	return true, true, nil
}
func (r *serviceRoleRepoStub) FindGlobalPermissionsByCodes(context.Context, repository.Transaction, []string) ([]entity.Permission, error) {
	return []entity.Permission{{ID: uuid.New(), Code: "can_view_user"}}, nil
}
func (r *serviceRoleRepoStub) CreateRole(_ context.Context, _ repository.Transaction, value entity.Role) (entity.Role, error) {
	value.ID = uuid.New()
	return value, nil
}
func (r *serviceRoleRepoStub) CreateRolePermission(_ context.Context, _ repository.Transaction, value entity.RolePermission) (entity.RolePermission, error) {
	return value, r.assignErr
}

func TestCreateRoleCommitsTransactionAndHashesToken(t *testing.T) {
	tx := &roleTransaction{}
	repo := &serviceRoleRepoStub{}
	_, err := NewRoleService(roleTransactionManager{tx: tx}, repo).CreateRole(context.Background(), "secret", "editor", "Редактор", []string{"can_view_user"})
	if err != nil {
		t.Fatal(err)
	}
	if !tx.committed {
		t.Fatal("transaction was not committed")
	}
	expected := sha256.Sum256([]byte("secret"))
	if repo.tokenHash != hex.EncodeToString(expected[:]) {
		t.Fatalf("token hash = %q", repo.tokenHash)
	}
}
func TestCreateRoleRollsBackWhenAssignmentFails(t *testing.T) {
	tx := &roleTransaction{}
	repo := &serviceRoleRepoStub{assignErr: errors.New("write failed")}
	_, err := NewRoleService(roleTransactionManager{tx: tx}, repo).CreateRole(context.Background(), "secret", "editor", "Редактор", []string{"can_view_user"})
	if err == nil {
		t.Fatal("expected error")
	}
	if tx.committed {
		t.Fatal("transaction must not be committed")
	}
	if !tx.rolledBack {
		t.Fatal("transaction was not rolled back")
	}
}
