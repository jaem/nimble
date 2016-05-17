package nim

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

// ResponseWriter is a wrapper around http.ResponseWriter that provides extra information about
// the response. It is recommended that middleware handlers use this construct to wrap a responsewriter
// if the functionality calls for it.
type Writer interface {
	http.ResponseWriter
	http.Flusher
	// Status returns the status code of the response or 0 if the response has not been written.
	Status() int
	// Written returns whether or not the ResponseWriter has been written.
	Written() bool
	// Size returns the size of the response body.
	Size() int
	// Before allows for a function to be called before the ResponseWriter has been written to. This is
	// useful for setting headers or any other operations that must happen before a response has been written.
	Before(func(Writer))
}

type beforeFunc func(Writer)

// NewResponseWriter creates a ResponseWriter that wraps an http.ResponseWriter
func newWriter(w http.ResponseWriter) Writer {
	return &writer{w, 0, 0, nil}
}

type writer struct {
	http.ResponseWriter
	status      int
	size        int
	beforeFuncs []beforeFunc
}

func (w *writer) WriteHeader(s int) {
	w.status = s
	w.callBefore()
	w.ResponseWriter.WriteHeader(s)
}

func (w *writer) Write(b []byte) (int, error) {
	if !w.Written() {
		// The status will be StatusOK if WriteHeader has not been called yet
		w.WriteHeader(http.StatusOK)
	}
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

func (w *writer) Status() int {
	return w.status
}

func (w *writer) Size() int {
	return w.size
}

func (w *writer) Written() bool {
	return w.status != 0
}

func (w *writer) Before(before func(Writer)) {
	w.beforeFuncs = append(w.beforeFuncs, before)
}

func (w *writer) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the ResponseWriter doesn't support the http.Hijack interface")
	}
	return hijacker.Hijack()
}

func (w *writer) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (w *writer) callBefore() {
	for i := len(w.beforeFuncs) - 1; i >= 0; i-- {
		w.beforeFuncs[i](w)
	}
}

func (w *writer) Flush() {
	flusher, ok := w.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}
