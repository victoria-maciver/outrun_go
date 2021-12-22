package main

import (
	"fmt"
	"net/http"

	"gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/dogs", handler).Methods("GET")

	http.ListenAndServe(":8080", r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}