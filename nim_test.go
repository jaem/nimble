package nim

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

func TestNimbleRun(t *testing.T) {
	// just test that Run doesn't bomb
	go New().Run(":3001")
}

func TestNimbleServeHTTP(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	n := New()
	n.UseHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "foo"
		next(w, r)
		result += "ban"
	})
	n.UseHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "bar"
		next(w, r)
		result += "baz"
	})
	n.UseHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "bat"
		w.WriteHeader(http.StatusBadRequest)
	})

	n.ServeHTTP(response, (*http.Request)(nil))

	expect(t, result, "foobarbatbazban")
	expect(t, response.Code, http.StatusBadRequest)
}

// Ensures that the middleware chain
// can correctly return all of its handlers.
func TestHandlers(t *testing.T) {
	response := httptest.NewRecorder()
	n := New()
	handles := n.handles
	expect(t, 0, len(handles))

	n.UseHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.WriteHeader(http.StatusOK)
	})

	// Expects the length of handlers to be exactly 1 
	// after adding exactly one handler to the middleware chain
	handles = n.handles
	expect(t, 1, len(handles))

	// Ensures that the first handler that is in sequence behaves
	// exactly the same as the one that was registered earlier
	handles[0](response, (*http.Request)(nil), nil)
	expect(t, response.Code, http.StatusOK)
}