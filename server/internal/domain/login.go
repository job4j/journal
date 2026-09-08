package domain

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
	"journal/server/internal/repository"
	"strings"
	"time"
)

type LoginRequest struct {
	Email    string
	Password string
}
type LoginResponse struct {
	User      repository.User
	Token     string
	ExpiresAt time.Time
}

func (d *AuthDomain) Login(ctx context.Context, tx repository.Transaction, request LoginRequest) (LoginResponse, error) {
	email := strings.TrimSpace(strings.ToLower(request.Email))
	if email == "" || request.Password == "" {
		return LoginResponse{}, ErrInvalidCredentials
	}
	user, err := d.repo.FindUserByEmail(ctx, tx, email)
	if errors.Is(err, repository.ErrNotFound) {
		return LoginResponse{}, ErrInvalidCredentials
	}
	if err != nil {
		return LoginResponse{}, fmt.Errorf("load user: %w", err)
	}
	valid, err := verifyArgon2ID(request.Password, user.PasswordHash)
	if err != nil || !valid {
		return LoginResponse{}, ErrInvalidCredentials
	}
	if user.Status != "active" {
		return LoginResponse{}, ErrUserBlocked
	}
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return LoginResponse{}, fmt.Errorf("generate session token: %w", err)
	}
	token := base64.RawURLEncoding.EncodeToString(tokenBytes)
	tokenHash := sha256.Sum256([]byte(token))
	expiresAt := d.now().UTC().Add(d.sessionLifetime)
	if err := d.repo.InsertSession(ctx, tx, repository.Session{ID: uuid.New(), UserID: user.ID, TokenHash: hex.EncodeToString(tokenHash[:]), ExpiresAt: expiresAt}); err != nil {
		return LoginResponse{}, fmt.Errorf("persist session: %w", err)
	}
	return LoginResponse{User: user, Token: token, ExpiresAt: expiresAt}, nil
}

func verifyArgon2ID(password, encoded string) (bool, error) {
	var memory uint32
	var iterations uint32
	var parallelism uint8
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" || parts[2] != "v=19" {
		return false, errors.New("invalid argon2id hash")
	}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism); err != nil {
		return false, err
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	expected, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	actual := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(expected)))
	return subtle.ConstantTimeCompare(actual, expected) == 1, nil
}
