package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func main() {
	proxyUrl, err := url.Parse("http://127.0.0.1:8888")
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}

	url := "http://ipinfo.io"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "curl/5.0")
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var a map[string]string
	err = json.Unmarshal(body, &a)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", a)
	fmt.Printf("%#v\n", a["ip"])
	fmt.Printf("%#v\n", res.Header)

}
