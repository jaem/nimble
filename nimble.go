package nimble

import (
	"log"
	"net/http"
	"os"

	"github.com/nimgo/nimble/interfaces"
	"github.com/nimgo/nimble/nimbleware"
)

// Nimble is a stack of Middleware Handlers that can be invoked as an http.Handler.
// The middleware stack is run in the sequence that they are added to the stack.
type Nimble struct {
	handlers   []HandlerFunc
	middleware middleware
}

// HandlerFunc is a linked-list handler interface that provides
// every middleware a forward reference to the next middleware in the stack.
type HandlerFunc func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

// Handler exposes an adapter to support specific middleware that uses this signature
// ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

// Each Middleware should yield to the next middleware in the chain by invoking the next http.HandlerFunc
type middleware struct {
	fn   HandlerFunc
	next *middleware
}

// Default returns a new Nimble instance with the default middleware already
// in the stack.
//
// Recovery - Panic Recovery Middleware
// Logger - Request/Response Logging
// Static - Static File Serving
func Default() *Nimble {
	return New().
		UseHandler(nimbleware.NewRecovery()).
		UseHandler(nimbleware.NewLogger()).
		UseHandler(nimbleware.NewStatic(http.Dir("static")))
}

// New returns a new Nimble instance with no middleware preconfigured.
func New() *Nimble {
	return &Nimble{}
}

// Run is a convenience function that runs the nimble stack as an HTTP
// server. The addr string takes the same format as http.ListenAndServe.
func (n *Nimble) Run(addr ...string) {
	l := log.New(os.Stdout, "[n.] ", 0)
	address := detectAddress(addr...)
	l.Printf("Server is listening on %s", address)
	l.Fatal(http.ListenAndServe(address, n))
}

// Nimble itself is a http.Handler. This allows it to used as a substack manager
func (n *Nimble) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, ok := w.(interfaces.Writer); ok { // handle substacks
		n.middleware.serve(w, r)
	} else {
		n.middleware.serve(newWriter(w), r)
	}
}

// Use adds a http.Handler onto the middleware stack.
func (n *Nimble) Use(handler http.Handler) *Nimble {
	return n.UseHandlerFunc(wrap(handler))
}

// UseFunc adds a http.HandlerFunc onto the middleware stack.
func (n *Nimble) UseFunc(handlerFunc http.HandlerFunc) *Nimble {
	return n.UseHandlerFunc(wrapHandlerFunc(handlerFunc))
}

// UseHandler adds a nimble.Handler onto the middleware stack.
func (n *Nimble) UseHandler(handler Handler) *Nimble {
	return n.UseHandlerFunc(handler.ServeHTTP)
}

// UseHandlerFunc adds a nimble.HandlerFunc function onto the middleware stack.
func (n *Nimble) UseHandlerFunc(handlerFunc HandlerFunc) *Nimble {
	if handlerFunc == nil {
		panic("handlerFunc cannot be nil")
	}
	n.handlers = append(n.handlers, handlerFunc)
	n.middleware = build(n.handlers)
	return n
}

// The next http.HandlerFunc is automatically called after the Handler is executed.
// If the Handler writes to the ResponseWriter, the next http.HandlerFunc should not be invoked.
func (m *middleware) serve(w http.ResponseWriter, r *http.Request) {
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
