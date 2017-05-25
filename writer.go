package nimble

import (
	"bufio"
	"fmt"
	"net"
	"net/http"

	"github.com/nimgo/nimble/interfaces"
)

// ResponseWriter is a wrapper around http.ResponseWriter that provides extra information about
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

type beforeFunc func(interfaces.Writer)

// newWriter creates a Writer that wraps an http.ResponseWriter
func newWriter(w http.ResponseWriter) interfaces.Writer {
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

func (w *writer) Before(before func(interfaces.Writer)) {
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
		flusher.Flush()
	}
}

// CloseNotify provides notifications when the HTTP connection terminates
func (wc *writerCloseNotifer) CloseNotify() <-chan bool {
	return wc.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
