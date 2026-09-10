package domain

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"testing"
	"time"

	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

type testTx struct{}

func (testTx) Commit(context.Context) error   { return nil }
func (testTx) Rollback(context.Context) error { return nil }

type authRepoStub struct {
	user         entity.User
	findErr      error
	inserted     *entity.Session
	passwordHash string
}

func (r *authRepoStub) FindUserByEmail(context.Context, repository.Transaction, string) (entity.User, error) {
	return r.user, r.findErr
}
func (r *authRepoStub) InsertSession(_ context.Context, _ repository.Transaction, session entity.Session) error {
	r.inserted = &session
	return nil
}
func (r *authRepoStub) FindUserBySessionTokenHash(context.Context, repository.Transaction, string, time.Time) (entity.User, error) {
	return r.user, r.findErr
}
func (r *authRepoStub) RevokeSession(context.Context, repository.Transaction, string, time.Time) error {
	return nil
}

func (r *authRepoStub) UpdatePasswordHash(
	_ context.Context,
	_ repository.Transaction,
	_ uuid.UUID,
	passwordHash string,
) error {
	r.passwordHash = passwordHash
	return nil
}
func TestLoginCreatesSessionForActiveUser(t *testing.T) {
	userID := uuid.New()
	repo := &authRepoStub{user: entity.User{ID: userID, Login: "user@example.com", PasswordHash: testHash("secret"), Status: "active"}}
	now := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	domain := NewAuthDomain(repo, time.Hour)
	domain.now = func() time.Time { return now }
	result, err := domain.Login(context.Background(), testTx{}, LoginRequest{Login: " USER@example.com ", Password: "secret"})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if result.Token == "" || repo.inserted == nil {
		t.Fatal("session was not created")
	}
	if repo.inserted.UserID != userID || !repo.inserted.ExpiresAt.Equal(now.Add(time.Hour)) {
		t.Fatalf("unexpected session: %+v", repo.inserted)
	}
}

func TestLoginRejectsInvalidPassword(t *testing.T) {
	repo := &authRepoStub{user: entity.User{PasswordHash: testHash("secret"), Status: "active"}}
	domain := NewAuthDomain(repo, time.Hour)
	_, err := domain.Login(context.Background(), testTx{}, LoginRequest{Login: "user@example.com", Password: "wrong"})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("error = %v", err)
	}
	if repo.inserted != nil {
		t.Fatal("session must not be inserted")
	}
}

func TestLoginRejectsBlockedUser(t *testing.T) {
	repo := &authRepoStub{user: entity.User{PasswordHash: testHash("secret"), Status: "blocked"}}
	domain := NewAuthDomain(repo, time.Hour)
	_, err := domain.Login(context.Background(), testTx{}, LoginRequest{Login: "user@example.com", Password: "secret"})
	if !errors.Is(err, ErrUserBlocked) {
		t.Fatalf("error = %v", err)
	}
	if repo.inserted != nil {
		t.Fatal("session must not be inserted")
	}
}

func testHash(password string) string {
	salt := []byte("0123456789abcdef")
	hash := argon2.IDKey([]byte(password), salt, 1, 64, 1, 32)
	return fmt.Sprintf("$argon2id$v=19$m=64,t=1,p=1$%s$%s", base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hash))
}

func TestChangeOwnPasswordVerifiesCurrentPassword(t *testing.T) {
	repo := &authRepoStub{user: entity.User{
		ID: uuid.New(), PasswordHash: testHash("secret"), Roles: []string{"teacher"},
	}}
	domain := NewAuthDomain(repo, time.Hour)

	if err := domain.ChangeOwnPassword(
		context.Background(), testTx{}, "hash", "secret", "x",
	); err != nil {
		t.Fatal(err)
	}
	valid, err := verifyArgon2ID("x", repo.passwordHash)
	if err != nil || !valid {
		t.Fatalf("new password was not stored: valid=%v error=%v", valid, err)
	}
}

func TestChangeOwnPasswordRejectsWrongCurrentPassword(t *testing.T) {
	repo := &authRepoStub{user: entity.User{
		ID: uuid.New(), PasswordHash: testHash("secret"), Roles: []string{"parent"},
	}}
	err := NewAuthDomain(repo, time.Hour).ChangeOwnPassword(
		context.Background(), testTx{}, "hash", "wrong", "x",
	)
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("error = %v", err)
	}
	if repo.passwordHash != "" {
		t.Fatal("password was changed")
	}
}

func TestChangeOwnPasswordRejectsStudent(t *testing.T) {
	repo := &authRepoStub{user: entity.User{
		ID: uuid.New(), PasswordHash: testHash("secret"), Roles: []string{"student"},
	}}
	err := NewAuthDomain(repo, time.Hour).ChangeOwnPassword(
		context.Background(), testTx{}, "hash", "secret", "x",
	)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("error = %v", err)
	}
}
