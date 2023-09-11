package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method == "GET" {

		fmt.Fprintf(w, "GET Method")

	} else if r.Method == "POST" {

		fmt.Fprintf(w, "POST Method\n%d\n", r.ContentLength)

	} else {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
	}

}

func main() {
	http.HandleFunc("/hello", helloHandler)

	fmt.Printf("Starting server at port 8083\n")
	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatal(err)
	}
}
