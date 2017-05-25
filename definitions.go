package nimble

import (
	"net/http"
	"os"
)

const (
	// DefaultAddress is used if no other is specified.
	defaultServerAddress = ":8080"
)

// detectAddress
func detectAddress(addr ...string) string {
	if len(addr) > 0 {
		return addr[0]
	}
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return defaultServerAddress
}

// Handler exposes an adapter to support specific middleware that uses this signature
// ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

// Nimble moves the stack by using a linked-list handler interface that provides
// every middleware a forward reference to the next middleware in the stack.
type HandlerFunc func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

// Each Middleware should yield to the next middleware in the chain by invoking the next http.HandlerFunc
type middleware struct {
	fn   HandlerFunc
	next *middleware
}

// The next http.HandlerFunc is automatically called after the Handler is executed.
// If the Handler writes to the ResponseWriter, the next http.HandlerFunc should not be invoked.
func (m middleware) serve(w http.ResponseWriter, r *http.Request) {
	m.fn(w, r, m.next.serve)
}

// Wrap converts a http.Handler into a nimble.HandlerFunc
func wrap(handler http.Handler) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handler.ServeHTTP(w, r)
		next(w, r)
	}
}

// wrapFunc converts a http.HandlerFunc into a nimble.HandlerFunc.
func wrapHandlerFunc(fn http.HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		fn(w, r)
		next(w, r)
	}
}

func build(handles []HandlerFunc) middleware {
	var next middleware

	if len(handles) == 0 {
		return emptyMiddleware()
	} else if len(handles) > 1 {
		next = build(handles[1:])
	} else {
		next = emptyMiddleware()
	}

	return middleware{handles[0], &next}
}

func emptyMiddleware() middleware {
	return middleware{
		func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) { /* do nothing */ },
		&middleware{},
	}
}
