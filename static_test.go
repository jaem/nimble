package nim

import (
	"bytes"
	"net/http"

	"testing"
	"net/http/httptest"

)

func TestStatic(t *testing.T) {
	response := httptest.NewRecorder()
	response.Body = new(bytes.Buffer)

	n := New()
	n.UseHandler(NewStatic(http.Dir(".")))

	req, err := http.NewRequest("GET", "http://localhost:3001/nim.go", nil)
	if err != nil {
		t.Error(err)
	}
	n.ServeHTTP(response, req)
	expect(t, response.Code, http.StatusOK)
	expect(t, response.Header().Get("Expires"), "")
	if response.Body.Len() == 0 {
		t.Errorf("Got empty body for GET request")
	}
}

func TestStaticHead(t *testing.T) {
	response := httptest.NewRecorder()
	response.Body = new(bytes.Buffer)

	n := New()
	n.UseHandler(NewStatic(http.Dir(".")))
	n.Use(http.NotFoundHandler())

	req, err := http.NewRequest("HEAD", "http://localhost:3001/nim.go", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	expect(t, response.Code, http.StatusOK)
	if response.Body.Len() != 0 {
		t.Errorf("Got non-empty body for HEAD request")
	}
}

func TestStaticAsPost(t *testing.T) {
	response := httptest.NewRecorder()

	n := New()
	n.UseHandler(NewStatic(http.Dir(".")))
	n.Use(http.NotFoundHandler())

	req, err := http.NewRequest("POST", "http://localhost:3001/nim.go", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	expect(t, response.Code, http.StatusNotFound)
}

func TestStaticBadDir(t *testing.T) {
	response := httptest.NewRecorder()

	n := Default()
	n.Use(http.NotFoundHandler())

	req, err := http.NewRequest("GET", "http://localhost:3001/nim.go", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	refute(t, response.Code, http.StatusOK)
}

func TestStaticOptionsServeIndex(t *testing.T) {
	response := httptest.NewRecorder()

	n := New()
	s := NewStatic(http.Dir("."))
	s.indexFile = "nim.go"
	n.UseHandler(s)

	req, err := http.NewRequest("GET", "http://localhost:3001/", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	expect(t, response.Code, http.StatusOK)
}

func TestStaticOptionsPrefix(t *testing.T) {
	response := httptest.NewRecorder()

	n := New()
	s := NewStatic(http.Dir("."))
	s.prefix = "/public"
	n.UseHandler(s)

	// Check file content behaviour
	req, err := http.NewRequest("GET", "http://localhost:3001/public/nim.go", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	expect(t, response.Code, http.StatusOK)
}
