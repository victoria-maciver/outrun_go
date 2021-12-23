package main

import (
	"fmt"
	"net/http"

	"gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

const (
	HelloWorld         string = "Hello World!"
	HelloWorldEndpoint string = "/hello-world"

	ListenPort string = ":8080"
	AssetsPath string = "/assets/"
)

func newRouter() *mux.Router {
	fmt.Println("Starting HTTP router")
	r := mux.NewRouter()
	h := NewHandler()
	r.HandleFunc(HelloWorldEndpoint, helloWorldHandler).Methods("GET") // testing endpoint
	r.HandleFunc("/dogs", h.listDogsHandler).Methods("GET")

	staticFileDirectory := http.Dir("." + AssetsPath)
	staticFileHandler := http.StripPrefix(AssetsPath, http.FileServer(staticFileDirectory))
	r.PathPrefix(AssetsPath).Handler(staticFileHandler).Methods("GET")

	return r
}

func main() {
	router := newRouter()
	fmt.Printf("Listening on port %v", ListenPort)
	http.ListenAndServe(ListenPort, router)
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, HelloWorld)
}
