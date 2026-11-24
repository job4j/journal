package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"testing"
)

func TestParentAttachmentsRequireBothPermissionsAndMembershipDates(t *testing.T) {
	for _, test := range []struct {
		name   string
		change func(*attachmentRepoStub)
		want   error
	}{
		{"enrollment day", func(*attachmentRepoStub) {}, nil},
		{"missing user permission", func(r *attachmentRepoStub) { r.denied = "can_view_user" }, ErrForbidden},
		{"missing journal permission", func(r *attachmentRepoStub) { r.denied = "can_view_journal" }, ErrForbidden},
		{"another class", func(r *attachmentRepoStub) { r.missingMembership = true }, ErrForbidden},
		{"before enrollment", func(r *attachmentRepoStub) { r.membership.EnrolledOn = r.lesson.LessonDate.AddDate(0, 0, 1) }, ErrForbidden},
		{"after leaving", func(r *attachmentRepoStub) {
			left := r.lesson.LessonDate.AddDate(0, 0, -1)
			r.membership.LeftOn = &left
		}, ErrForbidden},
		{"leaving day", func(r *attachmentRepoStub) { r.membership.LeftOn = &r.lesson.LessonDate }, nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := attachmentFixture()
			r.viewer.Roles = []string{"parent"}
			test.change(r)
			_, err := NewAttachmentDomain(r).ListLessonAttachments(context.Background(), testTx{}, "hash", ListLessonAttachmentsRequest{LessonID: r.lesson.ID, StudentID: uuid.New()})
			if !errors.Is(err, test.want) {
				t.Fatalf("got %v want %v", err, test.want)
			}
		})
	}
}
