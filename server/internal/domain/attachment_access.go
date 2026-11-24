package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"slices"
	"strings"
)

func (d *AttachmentDomain) authorize(ctx context.Context, tx repository.Transaction, hash string, lessonID, studentID uuid.UUID, write bool) (entity.User, error) {
	user, err := d.repo.FindActiveUserBySessionHash(ctx, tx, hash)
	if errors.Is(err, repository.ErrNotFound) {
		return user, ErrUnauthenticated
	}
	if err != nil {
		return user, err
	}
	lesson, err := d.repo.GetLesson(ctx, tx, lessonID)
	if errors.Is(err, repository.ErrNotFound) {
		return user, ErrLessonNotFound
	}
	if err != nil {
		return user, err
	}
	assignment, err := d.repo.GetClassSubject(ctx, tx, lesson.ClassSubjectID)
	if err != nil {
		return user, err
	}
	if studentID != uuid.Nil {
		if write || !slices.Contains(user.Roles, "parent") {
			return user, ErrForbidden
		}
		for _, code := range []string{"can_view_user", "can_view_journal"} {
			if err = AuthorizeObject(ctx, tx, d.repo, hash, code, studentID.String()); err != nil {
				return user, err
			}
		}
		student, err := d.repo.GetUser(ctx, tx, studentID)
		if err != nil {
			return user, err
		}
		if !slices.Contains(student.Roles, "student") {
			return user, ErrForbidden
		}
		membership, err := d.repo.GetClassStudent(ctx, tx, assignment.ClassID, studentID)
		if errors.Is(err, repository.ErrNotFound) {
			return user, ErrForbidden
		}
		if err != nil {
			return user, err
		}
		date := dayUTC(lesson.LessonDate)
		if date.Before(dayUTC(membership.EnrolledOn)) || (membership.LeftOn != nil && date.After(dayUTC(*membership.LeftOn))) {
			return user, ErrForbidden
		}
		return user, nil
	}
	if !slices.Contains(user.Roles, "admin") && (!slices.Contains(user.Roles, "teacher") || assignment.ResponsibleTeacherID != user.ID) {
		return user, ErrForbidden
	}
	code := "can_view_lesson"
	if write {
		code = "can_update_lesson"
	}
	if err = AuthorizeObject(ctx, tx, d.repo, hash, code, assignment.ID.String()); err != nil {
		return user, err
	}
	if write && strings.TrimSpace(lesson.Topic) == "" {
		return user, ErrInvalidLesson
	}
	return user, nil
}
