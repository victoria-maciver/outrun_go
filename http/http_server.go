package main

import (
	"fmt"
	"net/http"

	"gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

const helloWorld = "Hello World!"
const helloWorldEndpoint = "/hello-world"

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(helloWorldEndpoint, helloWorldHandler).Methods("GET")

	return r
}

func main() {
	r := newRouter()
	http.ListenAndServe(":8080", r)
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, helloWorld)
}
