package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
)

type ClassStudentView struct {
	Membership entity.ClassStudent
	Student    entity.User
}

func (d *ClassDomain) ListClassStudents(ctx context.Context, tx repository.Transaction, hash string, classID uuid.UUID) ([]ClassStudentView, error) {
	if err := AuthorizeObject(ctx, tx, d.repo, hash, "can_view_class", classID.String()); err != nil {
		return nil, err
	}
	if _, err := d.repo.GetClass(ctx, tx, classID); errors.Is(err, repository.ErrNotFound) {
		return nil, ErrClassNotFound
	} else if err != nil {
		return nil, fmt.Errorf("get class: %w", err)
	}
	memberships, err := d.repo.ListClassStudents(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list class students: %w", err)
	}
	users, err := d.repo.ListUsers(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	byID := map[uuid.UUID]entity.User{}
	for _, user := range users {
		byID[user.ID] = user
	}
	result := []ClassStudentView{}
	for _, membership := range memberships {
		if membership.ClassID == classID {
			if user, ok := byID[membership.UserID]; ok {
				result = append(result, ClassStudentView{Membership: membership, Student: user})
			}
		}
	}
	return result, nil
}
