package nim

import (
	"net/http"

	nctx "golang.org/x/net/context"

	"testing"
	"net/http/httptest"
)

func TestContext(t *testing.T) {
	recorder := httptest.NewRecorder()

	test_key := "key"
	test_val := 777777

	ctx := NewContext(nctx.TODO())

	n := New()
	// replace log for testing
	n.UseHandler(ctx)
	n.UseFunc(func(w http.ResponseWriter, r *http.Request) {
		c := GetContext(r)
		c = nctx.WithValue(c, test_key, test_val)
		SetContext(r, c)
	})
	n.UseFunc(func(w http.ResponseWriter, r *http.Request) {
		c := GetContext(r)
		if value, ok := c.Value(test_key).(int); ok {
			w.WriteHeader(value)
		}
	})

	req, err := http.NewRequest("GET", "http://localhost:3001/foobar", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(recorder, req)

	expect(t, recorder.Code, test_val)
}
