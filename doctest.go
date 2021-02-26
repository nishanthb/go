package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		mymux := http.NewServeMux()
		//mymux.HandleFunc("/", logger(funcOK))
		mymux.Handle("/", logger(http.HandlerFunc(funcOK)))
		s := &http.Server{
			Addr:    ":8080",
			Handler: mymux,
		}
		log.Fatal(s.ListenAndServe())

	} else {
		for _ = range [5]struct{}{} {
			resp, err := http.Get("http://" + os.Args[1] + ":8080")
			if err != nil {
				log.Printf("%v", err)
				continue
			}
			defer resp.Body.Close()
			b, _ := ioutil.ReadAll(resp.Body)
			log.Printf("resp is %v", string(b))
		}
		c := make(chan struct{}, 1)
		<-c
	}
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request: %v", r.URL.String())
		next.ServeHTTP(w, r)
	})
}
func funcOK(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// docker networking test
// create network
// docker network create mynet
// start server
// docker run -it -d --name=server --network=mynet server:latest
// start client
// docker run -it -d --name=client --network=mynet server:latest server
// check logs

// ./l.go:14:27: cannot use funcOK (type func(http.ResponseWriter, *http.Request)) as type http.Handler in argument to logger:
// func(http.ResponseWriter, *http.Request) does not implement http.Handler (missing ServeHTTP method)
// wrap it in http.HandlerFunc - http.HandlerFunc(funcOK)
