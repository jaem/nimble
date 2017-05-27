package nimbleware

import (
	"bytes"
	"log"
	"net/http"
	"reflect"

	"github.com/nimgo/nimble"

	"net/http/httptest"
	"testing"
)

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestRecovery(t *testing.T) {
	buff := bytes.NewBufferString("")
	rec := httptest.NewRecorder()

	recover := NewRecovery()
	recover.logger = log.New(buff, "[n.] ", 0)

	n := nimble.New()

	// replace log for testing
	n.WithHandler(recover)
	n.WithFunc(func(http.ResponseWriter, *http.Request) {
		panic("here is a panic!")
	})

	n.ServeHTTP(rec, (*http.Request)(nil))
	expect(t, rec.Code, http.StatusInternalServerError)
	refute(t, rec.Body.Len(), 0)
	refute(t, len(buff.String()), 0)
}
