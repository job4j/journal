package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
)

type ClassSubjectView struct {
	Assignment entity.ClassSubject
	Subject    entity.Subject
	Teacher    entity.User
}

func (d *ClassDomain) ListClassSubjects(ctx context.Context, tx repository.Transaction, hash string, classID uuid.UUID) ([]ClassSubjectView, error) {
	if err := AuthorizeObject(ctx, tx, d.repo, hash, "can_view_class", classID.String()); err != nil {
		return nil, err
	}
	if _, err := d.repo.GetClass(ctx, tx, classID); errors.Is(err, repository.ErrNotFound) {
		return nil, ErrClassNotFound
	} else if err != nil {
		return nil, fmt.Errorf("get class: %w", err)
	}
	assignments, err := d.repo.ListClassSubjects(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list class subjects: %w", err)
	}
	subjects, err := d.repo.ListSubjects(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list subjects: %w", err)
	}
	users, err := d.repo.ListUsers(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	subjectByID := make(map[uuid.UUID]entity.Subject, len(subjects))
	for _, subject := range subjects {
		subjectByID[subject.ID] = subject
	}
	userByID := make(map[uuid.UUID]entity.User, len(users))
	for _, user := range users {
		userByID[user.ID] = user
	}
	result := []ClassSubjectView{}
	for _, assignment := range assignments {
		if assignment.ClassID != classID {
			continue
		}
		_, allowed, err := d.repo.CheckSessionPermissionForValue(ctx, tx, hash, "can_view_class_subject", assignment.ID.String())
		if err != nil {
			return nil, fmt.Errorf("check class subject permission: %w", err)
		}
		if !allowed {
			continue
		}
		subject, subjectOK := subjectByID[assignment.SubjectID]
		teacher, teacherOK := userByID[assignment.ResponsibleTeacherID]
		if subjectOK && teacherOK {
			result = append(result, ClassSubjectView{Assignment: assignment, Subject: subject, Teacher: teacher})
		}
	}
	return result, nil
}
