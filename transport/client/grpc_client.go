package transport

import (
	"context"
	"io"
	"log"

	"github.com/angelus_reprobi/grpc_dog/pb"
	"google.golang.org/grpc"
)

type GrpcHandler struct {
	client pb.DogServiceClient
}

func New() *GrpcHandler {
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := pb.NewDogServiceClient(cc)
	return &GrpcHandler{c}
}

func (h *GrpcHandler) ListDogs() []pb.Dog {
	var dogs []pb.Dog

	stream, err := h.client.ListDog(
		context.Background(),
		&pb.ListDogRequest{},
	)
	if err != nil {
		log.Fatalf("error while calling ListDog RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		dogs = append(dogs, *res.GetDog())
	}

	return dogs
}
