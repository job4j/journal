package domain

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"testing"
	"time"

	"journal/server/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

type testTx struct{}

func (testTx) Commit(context.Context) error   { return nil }
func (testTx) Rollback(context.Context) error { return nil }

type authRepoStub struct {
	user     repository.User
	findErr  error
	inserted *repository.Session
}

func (r *authRepoStub) FindUserByEmail(context.Context, repository.Transaction, string) (repository.User, error) {
	return r.user, r.findErr
}
func (r *authRepoStub) InsertSession(_ context.Context, _ repository.Transaction, session repository.Session) error {
	r.inserted = &session
	return nil
}

func TestLoginCreatesSessionForActiveUser(t *testing.T) {
	userID := uuid.New()
	repo := &authRepoStub{user: repository.User{ID: userID, Email: "user@example.com", PasswordHash: testHash("secret"), Status: "active"}}
	now := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	result, err := Login(context.Background(), testTx{}, repo, " USER@example.com ", "secret", now, time.Hour)
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
	repo := &authRepoStub{user: repository.User{PasswordHash: testHash("secret"), Status: "active"}}
	_, err := Login(context.Background(), testTx{}, repo, "user@example.com", "wrong", time.Now(), time.Hour)
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("error = %v", err)
	}
	if repo.inserted != nil {
		t.Fatal("session must not be inserted")
	}
}

func TestLoginRejectsBlockedUser(t *testing.T) {
	repo := &authRepoStub{user: repository.User{PasswordHash: testHash("secret"), Status: "blocked"}}
	_, err := Login(context.Background(), testTx{}, repo, "user@example.com", "secret", time.Now(), time.Hour)
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
