package main

import (
    "fmt"
    "log"
    "net/http"
)

type countHandler struct {
    count int
}

func (c *countHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    c.count++
    fmt.Fprintf(w, "the count is %v", c.count)
}

// simple middleware
func logit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("url is %v", r.URL.String())
        next.ServeHTTP(w, r)
    })
}
func main() {

    n := new(countHandler)
    http.Handle("/count", logit(n))
    log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
