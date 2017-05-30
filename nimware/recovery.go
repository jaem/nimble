package nimware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
)

// NewRecovery returns a new instance of Recovery
func NewRecovery() *Recovery {
	return &Recovery{
		logger:     log.New(os.Stdout, "[nr.] ", 0),
		printStack: true,
		stackAll:   false,
		stackSize:  1024 * 8,
	}
}

// Recovery is a middleware that attempts to recover from panics and writes a 500 if there was one.
type Recovery struct {
	logger     ALogger
	printStack bool
	stackAll   bool
	stackSize  int
}

func (rec *Recovery) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {

			if w.Header().Get("Content-Type") == "" {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			}
			w.WriteHeader(http.StatusInternalServerError)

			stack := make([]byte, rec.stackSize)
			stack = stack[:runtime.Stack(stack, rec.stackAll)]

			f := "RECOVER: %s\n%s"
			rec.logger.Printf(f, err, stack)

			if rec.printStack {
				fmt.Fprintf(w, f, err, stack)
			}
		}
	}()

	next(w, r)
}
