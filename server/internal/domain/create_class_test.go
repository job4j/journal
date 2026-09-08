package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"testing"
)

func TestCreateClassValidatesAndCreates(t *testing.T) {
	item, err := NewClassDomain(classRepoStub{authenticated: true, allowed: true}).CreateClass(context.Background(), testTx{}, "hash", uuid.New(), " 7А ", 7)
	if err != nil {
		t.Fatal(err)
	}
	if item.Class.Name != "7А" || item.Class.ID == uuid.Nil {
		t.Fatalf("item = %+v", item)
	}
}
func TestCreateClassRejectsGrade(t *testing.T) {
	_, err := NewClassDomain(classRepoStub{authenticated: true, allowed: true}).CreateClass(context.Background(), testTx{}, "hash", uuid.New(), "12А", 12)
	if !errors.Is(err, ErrInvalidClass) {
		t.Fatalf("error = %v", err)
	}
}
