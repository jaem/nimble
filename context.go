package nimble

import (
	"net/http"

	nctx "golang.org/x/net/context"
	gctx "github.com/gorilla/context"
)

type key int
const contextkey key = 0

// Helper functions to get/set context
// Gorilla/mux is open to saving 100 params per request
// The idea is to save one context (with 100 params) per request
// Not sure how this affects performance, but it might reduce 
// the complexity of the sync.mutex in gorilla/context. 
//
// Nevertheless, this is a short-term implementation. Until 
// net/context arrives in http.Request.
func GetContext(r *http.Request) nctx.Context {
	if c, ok := gctx.GetOk(r, contextkey); ok {
		return c.(nctx.Context)
	}
	return nctx.TODO()
}

func SetContext(r *http.Request, c nctx.Context) {
	gctx.Set(r, contextkey, c)
}

// context is a middleware that provisions the context per request.
// context is not a context wrapper. It is a job that performs the
// task of context provisioning, generally at the start of the request. 
type context struct {
	baseContext nctx.Context
}

// NewContext returns a new context handler
func NewContext(c nctx.Context) *context {
	return &context{ baseContext: c }
}

// Performs the context provisioning as a middleware. Why middleware? 
// This allows for flexibility in usage. see nimble.DefaultWithContext()
func (c *context) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	gctx.Set(r, contextkey, c.baseContext)
	next(w, r)
}
