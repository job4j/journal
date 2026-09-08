package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"strings"
	"testing"
)

type userRepoStub struct {
	authenticated, allowed bool
	roles                  []entity.Role
	created                entity.User
	assigned               []entity.UserRole
}

func (r *userRepoStub) CheckSessionPermission(context.Context, repository.Transaction, string, string) (bool, bool, error) {
	return r.authenticated, r.allowed, nil
}
func (r *userRepoStub) CreateUser(_ context.Context, _ repository.Transaction, v entity.User) (entity.User, error) {
	v.ID = uuid.New()
	r.created = v
	return v, nil
}
func (r *userRepoStub) GetUser(context.Context, repository.Transaction, uuid.UUID) (entity.User, error) {
	return r.created, nil
}
func (r *userRepoStub) ListUsers(context.Context, repository.Transaction) ([]entity.User, error) {
	return []entity.User{r.created}, nil
}
func (r *userRepoStub) UpdateUser(_ context.Context, _ repository.Transaction, v entity.User) (entity.User, error) {
	r.created = v
	return v, nil
}
func (r *userRepoStub) DeleteUser(context.Context, repository.Transaction, uuid.UUID) error {
	return nil
}
func (r *userRepoStub) FindRolesByCodes(context.Context, repository.Transaction, []string) ([]entity.Role, error) {
	return r.roles, nil
}
func (r *userRepoStub) CreateUserRole(_ context.Context, _ repository.Transaction, v entity.UserRole) (entity.UserRole, error) {
	r.assigned = append(r.assigned, v)
	return v, nil
}
func (r *userRepoStub) ListUserRolesByUserID(context.Context, repository.Transaction, uuid.UUID) ([]entity.UserRole, error) {
	return r.assigned, nil
}
func (r *userRepoStub) DeleteUserRole(context.Context, repository.Transaction, uuid.UUID, uuid.UUID) error {
	return nil
}
func TestCreateUserHashesPasswordAndAssignsRoles(t *testing.T) {
	repo := &userRepoStub{authenticated: true, allowed: true, roles: []entity.Role{{ID: uuid.New(), Code: "teacher"}}}
	user, err := NewUserDomain(repo).CreateUser(context.Background(), nil, CreateUserRequest{SessionTokenHash: "hash", Input: UserInput{Email: " TEACHER@example.com ", Password: "password", FirstName: " Анна ", LastName: " Иванова ", Status: "active", Roles: []string{"teacher"}}})
	if err != nil {
		t.Fatal(err)
	}
	if user.Email != "teacher@example.com" || !strings.HasPrefix(repo.created.PasswordHash, "$argon2id$") {
		t.Fatalf("unexpected user: %+v", user)
	}
	if len(repo.assigned) != 1 {
		t.Fatalf("assignments = %d", len(repo.assigned))
	}
}
func TestCreateUserRequiresPermission(t *testing.T) {
	_, err := NewUserDomain(&userRepoStub{}).CreateUser(context.Background(), nil, CreateUserRequest{})
	if !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("error = %v", err)
	}
}
func TestCreateUserRejectsUnknownRole(t *testing.T) {
	repo := &userRepoStub{authenticated: true, allowed: true}
	_, err := NewUserDomain(repo).CreateUser(context.Background(), nil, CreateUserRequest{Input: UserInput{Email: "a@example.com", Password: "password", FirstName: "A", LastName: "B", Status: "active", Roles: []string{"missing"}}})
	if !errors.Is(err, ErrRoleNotFound) {
		t.Fatalf("error = %v", err)
	}
}
