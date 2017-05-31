package nimware

import (
	"bytes"
	"net/http"

	"net/http/httptest"
	"testing"

	"github.com/nimgo/nimble/nimble"
)

func TestStatic(t *testing.T) {
	rec := httptest.NewRecorder()
	rec.Body = new(bytes.Buffer)

	n := nimble.New()
	n.WithHandler(NewStatic(http.Dir(".")))

	req, err := http.NewRequest("GET", "http://localhost:3000/static_test.go", nil)
	if err != nil {
		t.Error(err)
	}
	n.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusOK)
	expect(t, rec.Header().Get("Expires"), "")
	if rec.Body.Len() == 0 {
		t.Errorf("Got empty body for GET request")
	}
}

func TestStaticHead(t *testing.T) {
	rec := httptest.NewRecorder()
	rec.Body = new(bytes.Buffer)

	n := nimble.New()
	n.WithHandler(NewStatic(http.Dir(".")))
	n.With(http.NotFoundHandler())

	req, err := http.NewRequest("HEAD", "http://localhost:3000/static_test.go", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusOK)
	if rec.Body.Len() != 0 {
		t.Errorf("Got non-empty body for HEAD request")
	}
}

func TestStaticAsPost(t *testing.T) {
	rec := httptest.NewRecorder()

	n := nimble.New()
	n.WithHandler(NewStatic(http.Dir(".")))
	n.With(http.NotFoundHandler())

	req, err := http.NewRequest("POST", "http://localhost:3000/static_test.go", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusNotFound)
}

func TestStaticBadDir(t *testing.T) {
	rec := httptest.NewRecorder()

	n := nimble.New().
		WithHandler(NewRecovery()).
		WithHandler(NewColorLogger()).
		WithHandler(NewStatic(http.Dir("static")))

	n.With(http.NotFoundHandler())

	req, err := http.NewRequest("GET", "http://localhost:3000/static_test.go", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(rec, req)
	refute(t, rec.Code, http.StatusOK)
}

func TestStaticOptionsServeIndex(t *testing.T) {
	rec := httptest.NewRecorder()

	n := nimble.New()
	s := NewStatic(http.Dir("."))
	s.indexFile = "nimble.go"
	n.WithHandler(s)

	req, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusOK)
}

func TestStaticOptionsPrefix(t *testing.T) {
	rec := httptest.NewRecorder()

	n := nimble.New()
	s := NewStatic(http.Dir("."))
	s.prefix = "/public"
	n.WithHandler(s)

	// Check file content behaviour
	req, err := http.NewRequest("GET", "http://localhost:3000/public/static_test.go", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusOK)
}
