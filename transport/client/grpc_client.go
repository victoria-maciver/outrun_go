package transport

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/angelus_reprobi/grpc_dog/pb"
	"google.golang.org/grpc"
)

type Client struct {
	client pb.DogServiceClient
}

var (
	srv  *Client
	once sync.Once
)

func NewClient() (*Client, error) {
	var (
		conn *grpc.ClientConn
		err  error
	)
	once.Do(func() {
		opts := grpc.WithInsecure()
		conn, err = grpc.Dial("0.0.0.0:50051", opts)
		client := pb.NewDogServiceClient(conn)

		if err == nil {
			srv = &Client{client}
		}
	})
	return srv, err
}

func (obj *Client) ListDogs() ([]pb.Dog, error) {
	var dogs []pb.Dog

	stream, grpcError := obj.client.ListDog(
		context.Background(),
		&pb.ListDogRequest{},
	)
	if grpcError != nil {
		return nil, fmt.Errorf("failed to call rpc ListDog : %s", grpcError.Error())
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to process received stream of rpc ListDog : %s", err.Error())
		}
		dogs = append(dogs, *res.GetDog())
	}

	return dogs, nil
}
