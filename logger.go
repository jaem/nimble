package nimble

import (
	"log"
	"os"
	"time"
	"net/http"
)

// Logger is a middleware that logs per request.
type Logger struct {
	*log.Logger
}

// NewLogger returns a new Logger instance
func NewLogger() *Logger {
	return &Logger{log.New(os.Stdout, "[n.] ", 0)}
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	l.Printf(CLR_G + "Request started %s %s" + CLR_N, r.Method, r.URL.Path)

	next(w, r)

	res := w.(Writer)
	l.Printf(CLR_C + "Request completed %v %s in %v" + CLR_N, res.Status(), http.StatusText(res.Status()), time.Since(start))
}

const CLR_0 = "\x1b[30;1m"
const CLR_R = "\x1b[31;1m"
const CLR_G = "\x1b[32;1m"
const CLR_Y = "\x1b[33;1m"
const CLR_B = "\x1b[34;1m"
const CLR_M = "\x1b[35;1m"
const CLR_C = "\x1b[36;1m"
const CLR_W = "\x1b[37;1m"
const CLR_N = "\x1b[0m"
