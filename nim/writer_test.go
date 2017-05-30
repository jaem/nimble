package nim

import (
	"bufio"
	"net"
	"net/http"
	"time"

	"net/http/httptest"
	"testing"
)

type closeNotifyingRecorder struct {
	*httptest.ResponseRecorder
	closed chan bool
}

func newCloseNotifyingRecorder() *closeNotifyingRecorder {
	return &closeNotifyingRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}

func (c *closeNotifyingRecorder) close() {
	c.closed <- true
}

func (c *closeNotifyingRecorder) CloseNotify() <-chan bool {
	return c.closed
}

type hijackableResponse struct {
	Hijacked bool
}

func newHijackableResponse() *hijackableResponse {
	return &hijackableResponse{}
}

func (h *hijackableResponse) Header() http.Header           { return nil }
func (h *hijackableResponse) Write(buf []byte) (int, error) { return 0, nil }
func (h *hijackableResponse) WriteHeader(code int)          {}
func (h *hijackableResponse) Flush()                        {}
func (h *hijackableResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h.Hijacked = true
	return nil, nil, nil
}

func TestResponseWriterWritingString(t *testing.T) {
	rec := httptest.NewRecorder()
	w := newWriter(rec)

	w.Write([]byte("Hello world"))

	expect(t, rec.Code, w.Status())
	expect(t, rec.Body.String(), "Hello world")
	expect(t, w.Status(), http.StatusOK)
	expect(t, w.Size(), 11)
	expect(t, w.Written(), true)
}

func TestResponseWriterWritingStrings(t *testing.T) {
	rec := httptest.NewRecorder()
	w := newWriter(rec)

	w.Write([]byte("Hello world"))
	w.Write([]byte("foo abc def ghi"))

	expect(t, rec.Code, w.Status())
	expect(t, rec.Body.String(), "Hello worldfoo abc def ghi")
	expect(t, w.Status(), http.StatusOK)
	expect(t, w.Size(), 26)
}

func TestResponseWriterWritingHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	w := newWriter(rec)

	w.WriteHeader(http.StatusNotFound)

	expect(t, rec.Code, w.Status())
	expect(t, rec.Body.String(), "")
	expect(t, w.Status(), http.StatusNotFound)
	expect(t, w.Size(), 0)
}

func TestResponseWriterBefore(t *testing.T) {
	rec := httptest.NewRecorder()
	w := newWriter(rec)
	result := ""

	w.Before(func(Writer) {
		result += "foo"
	})
	w.Before(func(Writer) {
		result += "bar"
	})

	w.WriteHeader(http.StatusNotFound)

	expect(t, rec.Code, w.Status())
	expect(t, rec.Body.String(), "")
	expect(t, w.Status(), http.StatusNotFound)
	expect(t, w.Size(), 0)
	expect(t, result, "barfoo")
}

func TestResponseWriterCloseNotify(t *testing.T) {
	rec := newCloseNotifyingRecorder()
	w := newWriter(rec)
	closed := false
	notifier := w.(http.CloseNotifier).CloseNotify()
	rec.close()
	select {
	case <-notifier:
		closed = true
	case <-time.After(time.Second):
	}
	expect(t, closed, true)
}

func TestResponseWriterNonCloseNotify(t *testing.T) {
	rec := httptest.NewRecorder()
	w := newWriter(rec)
	_, ok := w.(http.CloseNotifier)
	expect(t, ok, false)
}

func TestResponseWriterFlusher(t *testing.T) {
	rec := httptest.NewRecorder()
	w := newWriter(rec)

	_, ok := w.(http.Flusher)
	expect(t, ok, true)
}

func TestResponseWriter_Flush_marksWritten(t *testing.T) {
	rec := httptest.NewRecorder()
	w := newWriter(rec)

	w.Flush()
	expect(t, w.Status(), http.StatusOK)
	expect(t, w.Written(), true)
}

func TestResponseWriterHijack(t *testing.T) {
	hijackable := newHijackableResponse()
	w := newWriter(hijackable)
	hijacker, ok := w.(http.Hijacker)
	expect(t, ok, true)
	_, _, err := hijacker.Hijack()
	if err != nil {
		t.Error(err)
	}
	expect(t, hijackable.Hijacked, true)
}

func TestResponseWriteHijackNotOK(t *testing.T) {
	hijackable := new(http.ResponseWriter)
	w := newWriter(*hijackable)
	hijacker, ok := w.(http.Hijacker)
	expect(t, ok, true)
	_, _, err := hijacker.Hijack()

	refute(t, err, nil)
}
