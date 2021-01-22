package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type counter struct {
	m     sync.Mutex
	count int
}

// http.Handler function - requires ServeHTTP
func (c *counter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.m.Lock()
	defer c.m.Unlock()
	c.count++
	io.WriteString(w, "count is "+strconv.Itoa(c.count))
}

type Logger struct {
	l *log.Logger
}

func (l Logger) logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.l.Printf("url: %s", r.URL.String())
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	fh, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	l := log.New(fh, "logg: ", log.LstdFlags|log.LUTC|log.Lshortfile)
	logger := new(Logger)
	logger.l = l

	ctr := new(counter)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("404 found not")) })
	mux.Handle("/count", logger.logHandler(ctr))
	s := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}
	log.Fatal(s.ListenAndServe())
}
