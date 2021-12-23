package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/angelus_reprobi/grpc_dog/pb"
	"google.golang.org/grpc"
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

func newClient() pb.DogServiceClient {
	fmt.Println("Conecting to GRPC client")
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	return pb.NewDogServiceClient(cc)
}

func main() {
	router := newRouter()
	http.ListenAndServe(ListenPort, router)

	newClient()
	fmt.Println("HTTP server running")
	// client := newClient()
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, HelloWorld)
}
