package domain

import (
	"context"
	"errors"
	"journal/server/internal/repository"
	"testing"
)

type accessStub struct {
	authenticated, allowed bool
	err                    error
}

func (s accessStub) CheckSessionPermission(context.Context, repository.Transaction, string, string) (bool, bool, error) {
	return s.authenticated, s.allowed, s.err
}
func (s accessStub) CheckSessionPermissionForValue(context.Context, repository.Transaction, string, string, string) (bool, bool, error) {
	return s.authenticated, s.allowed, s.err
}

func TestAuthorize(t *testing.T) {
	for _, test := range []struct {
		name string
		stub accessStub
		want error
	}{
		{"allowed", accessStub{true, true, nil}, nil},
		{"unauthenticated", accessStub{}, ErrUnauthenticated},
		{"forbidden", accessStub{authenticated: true}, ErrForbidden},
	} {
		t.Run(test.name, func(t *testing.T) {
			if err := Authorize(context.Background(), testTx{}, test.stub, "hash", "permission"); !errors.Is(err, test.want) {
				t.Fatalf("error = %v, want %v", err, test.want)
			}
		})
	}
}

func TestAuthorizeObject(t *testing.T) {
	err := AuthorizeObject(context.Background(), testTx{}, accessStub{authenticated: true, allowed: true}, "hash", "can_view_user", "id")
	if err != nil {
		t.Fatal(err)
	}
}
