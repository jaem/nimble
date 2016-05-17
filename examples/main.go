package main

import (
"net/http"

"golang.org/x/net/context"

"github.com/gorilla/mux"
"github.com/nimgo/nim"
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
	router.PathPrefix("/p").Handler(nim.New().
		UseFunc(subMiddleware).
		Use(subrouter),
	)

	n := nim.DefaultWithContext(context.TODO())
	n.UseFunc(myMiddleware)
	n.Use(router)
	n.Run(":3000")
}

func aboutFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a lean, mean server."))
}

func myMiddleware(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("A middleware that always runs per http request.\n\n"))

	c := nim.GetContext(r)
	info := "ip = " + c.Value("ip").(string) + " port = " + c.Value("port").(string) + "\n\n"
	w.Write([]byte(info))

	c = context.WithValue(c, "key", "the Avengers")
	nim.SetContext(r, c)
}

func saysHi(who string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(who + " says, 'Hi y'all!'"))
	}
}

func subMiddleware(w http.ResponseWriter, r *http.Request) {
	c := nim.GetContext(r)
	if value, ok := c.Value("key").(string); ok {
		w.Write([]byte("SubMiddleware: Presenting to you " + value + "\n\n"))
	}
}
