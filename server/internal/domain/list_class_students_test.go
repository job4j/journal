package domain

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
	"testing"
)

func TestListClassStudentsReturnsHistory(t *testing.T) {
	classID, userID := uuid.New(), uuid.New()
	repo := classRepoStub{authenticated: true, allowed: true, classes: []entity.Class{{ID: classID}}, students: []entity.ClassStudent{{ClassID: classID, UserID: userID}}, users: []entity.User{{ID: userID, Name: "Иван", Roles: []string{"student"}}}}
	items, err := NewClassDomain(repo).ListClassStudents(context.Background(), testTx{}, "hash", classID)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Student.Name != "Иван" {
		t.Fatalf("items = %+v", items)
	}
}
