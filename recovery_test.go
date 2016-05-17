package nim

import (
	"bytes"
	"log"
	"net/http"

	"testing"
	"net/http/httptest"
)

func TestRecovery(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()

	rec := NewRecovery()
	rec.logger = log.New(buff, "[n.] ", 0)

	n := New()
	// replace log for testing
	n.UseHandler(rec)
	n.UseFunc(func(http.ResponseWriter, *http.Request) {
		panic("here is a panic!")
	})
	n.ServeHTTP(recorder, (*http.Request)(nil))
	expect(t, recorder.Code, http.StatusInternalServerError)
	refute(t, recorder.Body.Len(), 0)
	refute(t, len(buff.String()), 0)
}
