package api

import (
	"journal/server/internal/domain"
	"net/http/httptest"
	"testing"
)

func TestDeleteAttachmentHTTPContract(t *testing.T) {
	for _, test := range []struct {
		err    error
		status int
	}{{nil, 204}, {domain.ErrForbidden, 403}, {domain.ErrAttachmentNotFound, 404}, {domain.ErrAttachmentStorageUnavailable, 503}} {
		s := &attachmentServiceStub{err: test.err}
		res, err := attachmentApp(s).Test(httptest.NewRequest("DELETE", "/lesson-attachments/"+attachmentID.String(), nil))
		if err != nil {
			t.Fatal(err)
		}
		_ = res.Body.Close()
		if res.StatusCode != test.status || s.deleted != attachmentID {
			t.Fatalf("status=%d", res.StatusCode)
		}
	}
}
