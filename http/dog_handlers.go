package main

import (
	"encoding/json"
	"log"
	"net/http"

	transport "github.com/angelus_reprobi/grpc_dog/transport/client"
)

type HttpHandler struct {
	grpcHandler *transport.Client
}

func NewHandler() *HttpHandler {
	client, err := transport.NewClient()
	if err != nil {
		log.Fatalf("could not initialise grpc client : %v", err.Error())
	}
	return &HttpHandler{client}
}

func (h *HttpHandler) listDogsHandler(w http.ResponseWriter, r *http.Request) {
	dogList, err := h.grpcHandler.ListDogs()
	if err != nil {
		log.Fatalf("error calling ListDogs rpc : %v", err.Error())
	}

	dogListBytes, err := json.Marshal(dogList)
	if err != nil {
		log.Fatalf("Error processing doglist json : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(dogListBytes)
}
