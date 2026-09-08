package domain

import (
	"context"
	"errors"
	"fmt"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"strings"
)

func (d *SubjectDomain) ListSubjects(ctx context.Context, tx repository.Transaction, tokenHash string) ([]entity.Subject, error) {
	if err := Authorize(ctx, tx, d.repo, tokenHash, "can_view_subject"); err != nil {
		return nil, err
	}
	items, err := d.repo.ListSubjects(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list subjects: %w", err)
	}
	return items, nil
}
func (d *SubjectDomain) CreateSubject(ctx context.Context, tx repository.Transaction, tokenHash, code, name string) (entity.Subject, error) {
	if err := Authorize(ctx, tx, d.repo, tokenHash, "can_create_subject"); err != nil {
		return entity.Subject{}, err
	}
	code = strings.ToLower(strings.TrimSpace(code))
	name = strings.TrimSpace(name)
	if !roleCodePattern.MatchString(code) || name == "" || len([]rune(name)) > 100 {
		return entity.Subject{}, ErrInvalidSubject
	}
	item, err := d.repo.CreateSubject(ctx, tx, entity.Subject{Code: code, Name: name})
	if errors.Is(err, repository.ErrConflict) {
		return entity.Subject{}, ErrSubjectExists
	}
	if err != nil {
		return entity.Subject{}, fmt.Errorf("create subject: %w", err)
	}
	return item, nil
}
