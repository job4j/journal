package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func TestGrantParentStudentAccessIsIdempotent(t *testing.T) {
	parentID, studentID := uuid.New(), uuid.New()
	repo := &userRepoStub{authenticated: true, allowed: true, users: map[uuid.UUID]entity.User{parentID: {ID: parentID, Roles: []string{"parent"}}, studentID: {ID: studentID, Roles: []string{"student"}}}}
	domain := NewUserDomain(repo)
	for range 2 {
		if _, err := domain.GrantParentStudentAccess(context.Background(), testTx{}, "hash", parentID, studentID); err != nil {
			t.Fatal(err)
		}
	}
	if len(repo.permissions) != 2 || len(repo.userPermissions) != 2 {
		t.Fatalf("permissions=%d links=%d", len(repo.permissions), len(repo.userPermissions))
	}
}

func TestGrantParentStudentAccessValidatesRoles(t *testing.T) {
	parentID, studentID := uuid.New(), uuid.New()
	repo := &userRepoStub{authenticated: true, allowed: true, users: map[uuid.UUID]entity.User{parentID: {ID: parentID, Roles: []string{"teacher"}}, studentID: {ID: studentID, Roles: []string{"student"}}}}
	_, err := NewUserDomain(repo).GrantParentStudentAccess(context.Background(), testTx{}, "hash", parentID, studentID)
	if !errors.Is(err, ErrInvalidParentStudent) {
		t.Fatalf("error = %v", err)
	}
}
