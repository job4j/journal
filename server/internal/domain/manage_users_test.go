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
	users                  map[uuid.UUID]entity.User
	permissions            []entity.Permission
	userPermissions        []entity.UserPermission
	currentUser            entity.User
}

func (r *userRepoStub) CheckSessionPermission(context.Context, repository.Transaction, string, string) (bool, bool, error) {
	return r.authenticated, r.allowed, nil
}
func (r *userRepoStub) CreateUser(_ context.Context, _ repository.Transaction, v entity.User) (entity.User, error) {
	v.ID = uuid.New()
	r.created = v
	return v, nil
}

func (r *userRepoStub) GetUser(_ context.Context, _ repository.Transaction, id uuid.UUID) (entity.User, error) {
	if user, ok := r.users[id]; ok {
		return user, nil
	}
	if r.created.ID == uuid.Nil {
		return entity.User{}, repository.ErrNotFound
	}
	return r.created, nil
}
func (r *userRepoStub) ListUsers(context.Context, repository.Transaction) ([]entity.User, error) {
	if len(r.users) > 0 {
		result := make([]entity.User, 0, len(r.users))
		for _, user := range r.users {
			result = append(result, user)
		}
		return result, nil
	}
	return []entity.User{r.created}, nil
}
func (r *userRepoStub) FindActiveUserBySessionHash(context.Context, repository.Transaction, string) (entity.User, error) {
	if r.currentUser.ID == uuid.Nil {
		return entity.User{}, repository.ErrNotFound
	}
	return r.currentUser, nil
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
func (r *userRepoStub) EnsurePermission(_ context.Context, _ repository.Transaction, item entity.Permission) (entity.Permission, error) {
	for _, current := range r.permissions {
		if current.Code == item.Code && current.Value != nil && item.Value != nil && *current.Value == *item.Value {
			return current, nil
		}
	}
	item.ID = uuid.New()
	r.permissions = append(r.permissions, item)
	return item, nil
}
func (r *userRepoStub) EnsureUserPermission(_ context.Context, _ repository.Transaction, item entity.UserPermission) (entity.UserPermission, error) {
	for _, current := range r.userPermissions {
		if current == item {
			return current, nil
		}
	}
	r.userPermissions = append(r.userPermissions, item)
	return item, nil
}
func (r *userRepoStub) ListPermissions(context.Context, repository.Transaction) ([]entity.Permission, error) {
	return r.permissions, nil
}
func (r *userRepoStub) DeleteUserPermission(_ context.Context, _ repository.Transaction, userID, permissionID uuid.UUID) error {
	for i, item := range r.userPermissions {
		if item.UserID == userID && item.PermissionID == permissionID {
			r.userPermissions = append(r.userPermissions[:i], r.userPermissions[i+1:]...)
			return nil
		}
	}
	return repository.ErrNotFound
}
func (r *userRepoStub) ListUserPermissions(context.Context, repository.Transaction) ([]entity.UserPermission, error) {
	return r.userPermissions, nil
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
