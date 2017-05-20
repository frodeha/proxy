package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var client http.Client

func main() {

	http.HandleFunc("/s1/", createRedirectFunc("http://localhost:3000"))
	http.HandleFunc("/s2/", createRedirectFunc("http://localhost:4000"))
	log.Fatal(http.ListenAndServe(":2000", nil))
}

func createRedirectFunc(location string) func(w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse(location)
	if err != nil {
		panic(err)
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(url)
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Proxy to %s.\n", location)
		reverseProxy.ServeHTTP(w, r)
	}
}

func redirect(location string, w http.ResponseWriter, r *http.Request) {

	fmt.Printf("Redirecting to %s\n", location)
	// Build new url
	u, _ := url.Parse(location)
	u.Path = r.URL.Path
	q := u.Query()
	for k, v := range r.URL.Query() {
		for _, v1 := range v {
			q.Add(k, v1)
		}
	}
	u.RawQuery = q.Encode()

	// Create request
	req, err := http.NewRequest(r.Method, u.String(), r.Body)

	// Errors on malformed url
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for k, v := range r.Header {
		req.Header[k] = v
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	defer res.Body.Close()
	b, e := ioutil.ReadAll(res.Body)
	if e != nil {
		fmt.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		for k, v := range res.Header {
			for _, v1 := range v {
				w.Header().Set(k, v1)
			}
		}

		w.WriteHeader(res.StatusCode)
		w.Write(b)
	}
}
