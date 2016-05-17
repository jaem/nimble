package nim

import (
	"log"
	"os"
	"net/http"
	nctx "golang.org/x/net/context"
)

// New returns a new Nimble instance with no middleware preconfigured.
func New() *Nimble {
	return &Nimble{}
}

// Default returns a new Nimble instance with the default middleware already
// in the stack.
//
// Recovery - Panic Recovery Middleware
// Logger - Request/Response Logging
// Static - Static File Serving
// Context - Provision context
func Default() *Nimble {
	return New().
		UseHandler(NewRecovery()).
		UseHandler(NewLogger()).
		UseHandler(NewStatic(http.Dir("static")))
}

func DefaultWithContext(c nctx.Context) *Nimble {
	return New().
	UseHandler(NewRecovery()).
	UseHandler(NewLogger()).
	UseHandler(NewStatic(http.Dir("static"))).
	UseHandler(NewContext(c))
}

// Nimble is a stack of Middleware Handlers that can be invoked as an http.Handler.
// The middleware stack is run in the sequence that they are added to the stack.
type Nimble struct {
	middleware middleware
	handles   []Func
}

// Run is a convenience function that runs the nimble stack as an HTTP
// server. The addr string takes the same format as http.ListenAndServe.
func (n *Nimble) Run(addr string) {
	l := log.New(os.Stdout, "[n.] ", 0)
	l.Printf("Server listening on %s", addr)
	l.Fatal(http.ListenAndServe(addr, n))
}

// Nimble itself is a http.Handler. This allows it to used as a substack manager
func (n *Nimble) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, ok := w.(Writer); ok { // handle substacks
		n.middleware.serve(w, r)
	} else {
		n.middleware.serve(newWriter(w), r)
	}
}

// UseHandler adds a http.Handler onto the middleware stack.
func (n *Nimble) Use(handler http.Handler) *Nimble {
	return n.UseHandlerFunc(wrap(handler))
}

// UseHandlerFunc adds a http.HandlerFunc onto the middleware stack.
func (n *Nimble) UseFunc(handlerFunc http.HandlerFunc) *Nimble {
	return n.UseHandlerFunc(wrapFunc(handlerFunc))
}

// Use adds a nimble.Handler onto the middleware stack.
func (n *Nimble) UseHandler(handler Handler) *Nimble {
	return n.UseHandlerFunc(handler.ServeHTTP)
}

// UseFunc adds a nimble.Func function onto the middleware stack.
func (n *Nimble) UseHandlerFunc(fn Func) *Nimble {
	n.handles = append(n.handles, fn)
	n.middleware = build(n.handles)
	return n
}

// Handler exposes an adapter to support specific middleware that uses this signature
// ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

// Nimble moves the stack by using a linked-list handler interface that provides
// every middleware a forward reference to the next middleware in the stack.
type Func func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

// Each Middleware should yield to the next middleware in the chain by invoking the next http.HandlerFunc
type middleware struct {
	fn 		Func
	next  *middleware
}

// The next http.HandlerFunc is automatically called after the Handler is executed.
// If the Handler writes to the ResponseWriter, the next http.HandlerFunc should not be invoked.
func (m middleware) serve(w http.ResponseWriter, r *http.Request) {
	m.fn(w, r, m.next.serve)
}

// Wrap converts a http.Handler into a nimble.HandlerFunc
func wrap(handler http.Handler) Func {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handler.ServeHTTP(w, r)
		next(w, r)
	}
}

// wrapFunc converts a http.HandlerFunc into a nimble.HandlerFunc.
func wrapFunc(fn http.HandlerFunc) Func {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		fn(w, r)
		next(w, r)
	}
}

func build(handles []Func) middleware {
	var next middleware

	if len(handles) == 0 {
		return empty()
	} else if len(handles) > 1 {
		next = build(handles[1:])
	} else {
		next = empty()
	}

	return middleware{ handles[0], &next }
}

func empty() middleware {
	return middleware{
		func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) { /* do nothing */ },
		&middleware{},
	}
}
