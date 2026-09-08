package domain

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"net/mail"
	"sort"
	"strings"
)

type UserDomain struct {
	repo repository.UserManagementRepository
}

func NewUserDomain(repo repository.UserManagementRepository) *UserDomain {
	return &UserDomain{repo: repo}
}

type UserInput struct {
	Email, Password, FirstName, LastName, Status string
	Roles                                        []string
}
type CreateUserRequest struct {
	SessionTokenHash string
	Input            UserInput
}
type UpdateUserRequest struct {
	SessionTokenHash string
	ID               uuid.UUID
	Input            UserInput
}
type ListUsersRequest struct {
	SessionTokenHash, Role string
	Limit, Offset          int
}
type GetUserRequest struct {
	SessionTokenHash string
	ID               uuid.UUID
}
type DeleteUserRequest struct {
	SessionTokenHash string
	ID               uuid.UUID
}
type ListUsersResult struct {
	Items []entity.User
	Total int
}

func (d *UserDomain) authorize(ctx context.Context, tx repository.Transaction, hash, permission string) error {
	authenticated, allowed, err := d.repo.CheckSessionPermission(ctx, tx, hash, permission)
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
func normalizeUser(input UserInput, passwordRequired bool) (UserInput, error) {
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.FirstName = strings.TrimSpace(input.FirstName)
	input.LastName = strings.TrimSpace(input.LastName)
	address, err := mail.ParseAddress(input.Email)
	if err != nil || address.Address != input.Email || input.FirstName == "" || input.LastName == "" || len([]rune(input.FirstName)) > 100 || len([]rune(input.LastName)) > 100 {
		return UserInput{}, ErrInvalidUser
	}
	if input.Status != "active" && input.Status != "blocked" {
		return UserInput{}, ErrInvalidUser
	}
	if passwordRequired && len([]rune(input.Password)) < 8 {
		return UserInput{}, ErrInvalidUser
	}
	if input.Password != "" && len([]rune(input.Password)) < 8 {
		return UserInput{}, ErrInvalidUser
	}
	seen := map[string]struct{}{}
	roles := make([]string, 0, len(input.Roles))
	for _, raw := range input.Roles {
		code := strings.ToLower(strings.TrimSpace(raw))
		if !roleCodePattern.MatchString(code) {
			return UserInput{}, ErrInvalidUser
		}
		if _, ok := seen[code]; ok {
			return UserInput{}, ErrInvalidUser
		}
		seen[code] = struct{}{}
		roles = append(roles, code)
	}
	if len(roles) == 0 {
		return UserInput{}, ErrInvalidUser
	}
	sort.Strings(roles)
	input.Roles = roles
	return input, nil
}
func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, 3, 64*1024, 2, 32)
	return fmt.Sprintf("$argon2id$v=19$m=65536,t=3,p=2$%s$%s", base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hash)), nil
}
func (d *UserDomain) resolveRoles(ctx context.Context, tx repository.Transaction, codes []string) ([]entity.Role, error) {
	roles, err := d.repo.FindRolesByCodes(ctx, tx, codes)
	if err != nil {
		return nil, fmt.Errorf("load roles: %w", err)
	}
	if len(roles) != len(codes) {
		return nil, ErrRoleNotFound
	}
	return roles, nil
}
func (d *UserDomain) CreateUser(ctx context.Context, tx repository.Transaction, request CreateUserRequest) (entity.User, error) {
	if err := d.authorize(ctx, tx, request.SessionTokenHash, "can_create_user"); err != nil {
		return entity.User{}, err
	}
	input, err := normalizeUser(request.Input, true)
	if err != nil {
		return entity.User{}, err
	}
	roles, err := d.resolveRoles(ctx, tx, input.Roles)
	if err != nil {
		return entity.User{}, err
	}
	passwordHash, err := hashPassword(input.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("hash password: %w", err)
	}
	user, err := d.repo.CreateUser(ctx, tx, entity.User{Email: input.Email, PasswordHash: passwordHash, FirstName: input.FirstName, LastName: input.LastName, Status: input.Status})
	if errors.Is(err, repository.ErrConflict) {
		return entity.User{}, ErrUserExists
	}
	if err != nil {
		return entity.User{}, fmt.Errorf("create user: %w", err)
	}
	for _, role := range roles {
		if _, err = d.repo.CreateUserRole(ctx, tx, entity.UserRole{UserID: user.ID, RoleID: role.ID}); err != nil {
			return entity.User{}, fmt.Errorf("assign role: %w", err)
		}
	}
	user.Roles = input.Roles
	return user, nil
}
func (d *UserDomain) ListUsers(ctx context.Context, tx repository.Transaction, request ListUsersRequest) (ListUsersResult, error) {
	if err := d.authorize(ctx, tx, request.SessionTokenHash, "can_view_user"); err != nil {
		return ListUsersResult{}, err
	}
	users, err := d.repo.ListUsers(ctx, tx)
	if err != nil {
		return ListUsersResult{}, fmt.Errorf("list users: %w", err)
	}
	filtered := make([]entity.User, 0, len(users))
	for _, user := range users {
		if request.Role == "" || containsCode(user.Roles, request.Role) {
			filtered = append(filtered, user)
		}
	}
	total := len(filtered)
	start := request.Offset
	if start > total {
		start = total
	}
	end := start + request.Limit
	if end > total {
		end = total
	}
	return ListUsersResult{Items: filtered[start:end], Total: total}, nil
}
func (d *UserDomain) GetUser(ctx context.Context, tx repository.Transaction, request GetUserRequest) (entity.User, error) {
	authenticated, allowed, err := d.repo.CheckSessionPermissionForValue(ctx, tx, request.SessionTokenHash, "can_view_user", request.ID.String())
	if err != nil {
		return entity.User{}, fmt.Errorf("check access: %w", err)
	}
	if !authenticated {
		return entity.User{}, ErrUnauthenticated
	}
	if !allowed {
		return entity.User{}, ErrForbidden
	}
	user, err := d.repo.GetUser(ctx, tx, request.ID)
	if errors.Is(err, repository.ErrNotFound) {
		return entity.User{}, ErrUserNotFound
	}
	if err != nil {
		return entity.User{}, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}
func (d *UserDomain) UpdateUser(ctx context.Context, tx repository.Transaction, request UpdateUserRequest) (entity.User, error) {
	if err := d.authorize(ctx, tx, request.SessionTokenHash, "can_update_user"); err != nil {
		return entity.User{}, err
	}
	input, err := normalizeUser(request.Input, false)
	if err != nil {
		return entity.User{}, err
	}
	current, err := d.repo.GetUser(ctx, tx, request.ID)
	if errors.Is(err, repository.ErrNotFound) {
		return entity.User{}, ErrUserNotFound
	}
	if err != nil {
		return entity.User{}, fmt.Errorf("get user: %w", err)
	}
	roles, err := d.resolveRoles(ctx, tx, input.Roles)
	if err != nil {
		return entity.User{}, err
	}
	passwordHash := current.PasswordHash
	if input.Password != "" {
		passwordHash, err = hashPassword(input.Password)
		if err != nil {
			return entity.User{}, fmt.Errorf("hash password: %w", err)
		}
	}
	user, err := d.repo.UpdateUser(ctx, tx, entity.User{ID: current.ID, Email: input.Email, PasswordHash: passwordHash, FirstName: input.FirstName, LastName: input.LastName, Status: input.Status})
	if errors.Is(err, repository.ErrConflict) {
		return entity.User{}, ErrUserExists
	}
	if err != nil {
		return entity.User{}, fmt.Errorf("update user: %w", err)
	}
	links, err := d.repo.ListUserRolesByUserID(ctx, tx, current.ID)
	if err != nil {
		return entity.User{}, fmt.Errorf("list user roles: %w", err)
	}
	for _, link := range links {
		if err = d.repo.DeleteUserRole(ctx, tx, link.UserID, link.RoleID); err != nil {
			return entity.User{}, fmt.Errorf("remove role: %w", err)
		}
	}
	for _, role := range roles {
		if _, err = d.repo.CreateUserRole(ctx, tx, entity.UserRole{UserID: user.ID, RoleID: role.ID}); err != nil {
			return entity.User{}, fmt.Errorf("assign role: %w", err)
		}
	}
	user.Roles = input.Roles
	return user, nil
}
func (d *UserDomain) DeleteUser(ctx context.Context, tx repository.Transaction, request DeleteUserRequest) error {
	if err := d.authorize(ctx, tx, request.SessionTokenHash, "can_update_user"); err != nil {
		return err
	}
	err := d.repo.DeleteUser(ctx, tx, request.ID)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrUserNotFound
	}
	if errors.Is(err, repository.ErrConflict) || errors.Is(err, repository.ErrReference) {
		return ErrUserInUse
	}
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
