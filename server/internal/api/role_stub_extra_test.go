package api

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/domain"
)

func (s *roleServiceStub) ListRoles(context.Context, string) (domain.ListRolesResponse, error) {
	return domain.ListRolesResponse{}, s.err
}
func (s *roleServiceStub) ListPermissions(context.Context, string) (domain.ListPermissionsResponse, error) {
	return domain.ListPermissionsResponse{}, s.err
}
func (s *roleServiceStub) UpdateRole(context.Context, string, uuid.UUID, string, string, []string) (domain.RoleView, error) {
	return domain.RoleView{}, s.err
}
func (s *roleServiceStub) DeleteRole(context.Context, string, uuid.UUID) error { return s.err }
