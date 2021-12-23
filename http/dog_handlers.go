package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	transport "github.com/angelus_reprobi/grpc_dog/transport/client"
)

type HttpHandler struct {
	grpcHandler *transport.GrpcHandler
}

func NewHandler() *HttpHandler {
	return &HttpHandler{transport.New()}
}

func (h *HttpHandler) listDogsHandler(w http.ResponseWriter, r *http.Request) {
	dogList := h.grpcHandler.ListDogs()

	dogListBytes, err := json.Marshal(dogList)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(dogListBytes)
}
