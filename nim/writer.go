package nim

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

// Writer is the interface response wrapper that provides extra information about
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

// ResponseWriter is a light wrapper around http.ResponseWriter that provides extra information about
// the response. It is recommended that middleware handlers use this construct to wrap a responsewriter
// if the functionality calls for it.
type writer struct {
	http.ResponseWriter
	status      int
	size        int
	beforeFuncs []beforeFunc
}

type writerCloseNotifer struct {
	*writer
}

type beforeFunc func(Writer)

// newWriter creates a Writer that wraps an http.ResponseWriter
func newWriter(w http.ResponseWriter) Writer {
	nw := &writer{
		ResponseWriter: w,
	}

	// provide closenotifier calls only if the writer implements it
	if _, ok := w.(http.CloseNotifier); ok {
		return &writerCloseNotifer{nw}
	}

	return nw
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

func (w *writer) callBefore() {
	for i := len(w.beforeFuncs) - 1; i >= 0; i-- {
		w.beforeFuncs[i](w)
	}
}

func (w *writer) Flush() {
	flusher, ok := w.ResponseWriter.(http.Flusher)
	if ok {
		if !w.Written() {
			// The status will be StatusOK if WriteHeader has not been called yet
			w.WriteHeader(http.StatusOK)
		}
		flusher.Flush()
	}
}

// CloseNotify provides notifications when the HTTP connection terminates
func (wcn *writerCloseNotifer) CloseNotify() <-chan bool {
	return wcn.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
