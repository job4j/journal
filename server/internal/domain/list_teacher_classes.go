package domain

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"journal/server/internal/repository"
)

func (d *ClassDomain) currentTeacher(ctx context.Context, tx repository.Transaction, hash string) (uuid.UUID, error) {
	user, err := d.repo.FindActiveUserBySessionHash(ctx, tx, hash)
	if errors.Is(err, repository.ErrNotFound) {
		return uuid.Nil, ErrUnauthenticated
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("find session user: %w", err)
	}
	if !slices.Contains(user.Roles, "teacher") {
		return uuid.Nil, ErrForbidden
	}
	return user.ID, nil
}

func (d *ClassDomain) ListTeacherClasses(ctx context.Context, tx repository.Transaction, hash string) ([]ClassView, error) {
	teacherID, err := d.currentTeacher(ctx, tx, hash)
	if err != nil {
		return nil, err
	}
	assignments, err := d.repo.ListClassSubjects(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list class subjects: %w", err)
	}
	allowedClasses := map[uuid.UUID]struct{}{}
	for _, item := range assignments {
		if item.ResponsibleTeacherID != teacherID {
			continue
		}
		_, allowed, checkErr := d.repo.CheckSessionPermissionForValue(ctx, tx, hash, "can_view_class_subject", item.ID.String())
		if checkErr != nil {
			return nil, fmt.Errorf("check assignment permission: %w", checkErr)
		}
		if allowed {
			allowedClasses[item.ClassID] = struct{}{}
		}
	}
	classes, err := d.repo.ListClasses(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list classes: %w", err)
	}
	students, err := d.repo.ListClassStudents(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list class students: %w", err)
	}
	counts := map[uuid.UUID]int{}
	for _, item := range students {
		if item.LeftOn == nil {
			counts[item.ClassID]++
		}
	}
	result := []ClassView{}
	for _, item := range classes {
		if _, ok := allowedClasses[item.ID]; ok {
			result = append(result, ClassView{Class: item, StudentCount: counts[item.ID]})
		}
	}
	return result, nil
}

func (d *ClassDomain) ListTeacherClassSubjects(ctx context.Context, tx repository.Transaction, hash string, classID uuid.UUID) ([]ClassSubjectView, error) {
	teacherID, err := d.currentTeacher(ctx, tx, hash)
	if err != nil {
		return nil, err
	}
	if _, err = d.repo.GetClass(ctx, tx, classID); errors.Is(err, repository.ErrNotFound) {
		return nil, ErrClassNotFound
	} else if err != nil {
		return nil, fmt.Errorf("get class: %w", err)
	}
	if err = AuthorizeObject(ctx, tx, d.repo, hash, "can_view_class", classID.String()); err != nil {
		return nil, err
	}
	items, err := d.ListClassSubjects(ctx, tx, hash, classID)
	if err != nil {
		return nil, err
	}
	result := []ClassSubjectView{}
	for _, item := range items {
		if item.Assignment.ResponsibleTeacherID == teacherID {
			result = append(result, item)
		}
	}
	return result, nil
}
