// Package nimble is a lightweight middleware manager in Golang that is based on net/http standards.
//
// For a full guide visit http://github.com/jaem/nimble
//
//  package main
//
//  import (
//    "net/http"
//    "fmt"
//    "github.com/jaem/nimble"
//  )
//
//  func main() {
//    mux := http.NewServeMux()
//    mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
//      fmt.Fprintf(w, "Welcome to nimble!")
//    })
//
//    n := nimble.Default()
//    n.Use(mux)
//    n.Run(":8000")
//  }
package nim
