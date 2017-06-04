package nimble

import (
	"net/http"
	"net/http/httptest"
	"reflect"
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

func TestNimbleServeHTTP(t *testing.T) {
	rec := httptest.NewRecorder()

	result := ""

	n := New()
	n.WithHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "_1bef"
		next(w, r)
		result += "_1aft"
	})
	n.WithHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "_2bef"
		next(w, r)
		result += "_2aft"
	})
	n.WithHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "_3here"
		w.WriteHeader(http.StatusBadRequest)
	})

	n.ServeHTTP(rec, (*http.Request)(nil))

	expect(t, result, "_1bef_2bef_3here_2aft_1aft")
	expect(t, rec.Code, http.StatusBadRequest)
}

// Ensures that the middleware chain
// can correctly return all of its handlers.
func TestUseHandlerFunc(t *testing.T) {
	rec := httptest.NewRecorder()
	n := New()
	handles := n.handlers
	expect(t, 0, len(handles))

	n.WithHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.WriteHeader(http.StatusOK)
	})

	// Expects the length of handlers to be exactly 1
	// after adding exactly one handler to the middleware chain
	handles = n.handlers
	expect(t, 1, len(handles))

	// Ensures that the first handler that is in sequence behaves
	// exactly the same as the one that was registered earlier
	handles[0](rec, (*http.Request)(nil), nil)
	expect(t, rec.Code, http.StatusOK)
}

func TestNimbleUseNil(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("Expected nimble.Use(nil) to panic, but it did not")
		}
	}()

	n := New()
	n.With(nil)
}
