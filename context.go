package nimble

import (
	"net/http"

	nctx "golang.org/x/net/context"
	gctx "github.com/gorilla/context"
	"net"
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

	ip, port := getIPaddress(w, r)

	ctx := nctx.WithValue(c.baseContext, "ip", ip)
	ctx  = nctx.WithValue(ctx, "port", port)

	gctx.Set(r, contextkey, ctx)

	next(w, r)
}


// http://stackoverflow.com/questions/27234861/correct-way-of-getting-clients-ip-addresses-from-http-request-golang
func getIPaddress(w http.ResponseWriter, r *http.Request) (string, string) {
	// This will only be defined when site is accessed via non-anonymous proxy
	// and takes precedence over RemoteAddr
	// Header.Get is case-insensitive
	if ipProxy := r.Header.Get("x-forwarded-for"); len(ipProxy) > 0 {
		return ipProxy, ""
	}

	ip, port, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return ip, port // r.RemoteAddr contains IP address
	}

	userIP := net.ParseIP(ip)
	if userIP != nil {
		return userIP.String(), "" // too bad??
	}

	return "", ""
}
