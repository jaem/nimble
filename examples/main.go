package main

import (
"net/http"

"golang.org/x/net/context"

"github.com/gorilla/mux"
"github.com/jaem/nimble"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Nimble!"))
	})
	router.HandleFunc("/about", aboutFunc)

	subrouter := mux.NewRouter()
	subrouter.HandleFunc("/p/iron_man", saysHi("Iron Man"))
	subrouter.HandleFunc("/p/captain_america", saysHi("Captain America"))
	router.PathPrefix("/p").Handler(nimble.New().
		UseFunc(subMiddleware).
		Use(subrouter),
	)

	n := nimble.Default()
	n.UseFunc(myMiddleware)
	n.Use(router)
	n.Run(":3000")
}

func aboutFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a lean, mean server."))
}

func myMiddleware(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("A middleware that always runs per http request.\n\n"))

	c := nimble.GetContext(r)
	c = context.WithValue(c, "key", "the Avengers")
	nimble.SetContext(r, c)
}

func saysHi(who string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(who + " says, 'Hi y'all!'"))
	}
}

func subMiddleware(w http.ResponseWriter, r *http.Request) {
	c := nimble.GetContext(r)
	if value, ok := c.Value("key").(string); ok {
		w.Write([]byte("SubMiddleware: Presenting to you " + value + "\n\n"))
	}
}
