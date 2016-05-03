package nimble

import (
	"net/http"

	"golang.org/x/net/context"
	gorilla "github.com/gorilla/context"
)

type key int
const contextkey key = 0

func GetContext(r *http.Request) context.Context {
	if c, ok := gorilla.GetOk(r, contextkey); ok {
		return c.(context.Context)
	}
	return context.TODO()
}

func SetContext(r *http.Request, c context.Context) {
	gorilla.Set(r, contextkey, c)
}

// Logger is a middleware that logs per request.
type ncontext struct {
	baseContext context.Context
}

// NewLogger returns a new Logger instance
func NewContext(c context.Context) *ncontext {
	return &ncontext{ baseContext: c }
}

func (c *ncontext) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	gorilla.Set(r, contextkey, c.baseContext)
	next(w, r)
}