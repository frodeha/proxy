package main

import "net/http"
import "log"

func main() {

	http.HandleFunc("/", redirect)
	log.Fatal(http.ListenAndServe(":2000", nil))
}

func redirect(w http.ResponseWriter, r *http.Request) {

}
