package nim

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"net/http"
)

// NewRecovery returns a new instance of Recovery
func NewRecovery() *recovery {
	return &recovery{
		logger:     log.New(os.Stdout, "[nim.] ", 0),
		printStack: true,
		stackAll:   false,
		stackSize:  1024 * 8,
	}
}

// Recovery is a middleware that attempts to recover from panics and writes a 500 if there was one.
type recovery struct {
	logger     *log.Logger
	printStack bool
	stackAll   bool
	stackSize  int
}

func (rec *recovery) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
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
