package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/angelus_reprobi/grpc_dog/src/pb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"gopkg.in/mgo.v2/bson"
)

var collection *mongo.Collection

type server struct {
	pb.DogServiceServer
}

type dogData struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	RegName  string             `bson:"reg_name"`
	CallName string             `bson:"call_name"`
	Gender   string             `bson:"gender"`
}

func (*server) CreateDog(ctx context.Context, req *pb.CreateDogRequest) (*pb.CreateDogResponse, error) {
	dog := req.GetDog()

	data := dogData{
		RegName:  dog.GetRegName(),
		CallName: dog.GetCallName(),
		Gender:   dog.GetGender(),
	}

	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to OID"),
		)
	}

	return &pb.CreateDogResponse{
		Dog: &pb.Dog{
			Id:       oid.Hex(),
			RegName:  dog.GetRegName(),
			CallName: dog.GetCallName(),
			Gender:   dog.GetGender(),
		},
	}, nil
}

func (*server) GetDog(ctx context.Context, req *pb.GetDogRequest) (*pb.GetDogResponse, error) {
	dogId := req.GetDogId()
	oid, err := primitive.ObjectIDFromHex(dogId)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	data := &dogData{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find dog with specified ID: %v", err),
		)
	}

	return &pb.GetDogResponse{
		Dog: dataToDogPb(data),
	}, nil
}

func (*server) UpdateDog(ctx context.Context, req *pb.UpdateDogRequest) (*pb.UpdateDogResponse, error) {
	dog := req.GetDog()
	oid, err := primitive.ObjectIDFromHex(dog.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	data := &dogData{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find dog with specified ID: %v", err),
		)
	}

	data.RegName = dog.GetRegName()
	data.CallName = dog.GetCallName()
	data.Gender = dog.GetGender()

	_, updateError := collection.ReplaceOne(context.Background(), filter, data)
	if updateError != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot update object in MongoDB: %v", updateError),
		)
	}

	return &pb.UpdateDogResponse{
		Dog: dataToDogPb(data),
	}, nil
}

func (*server) ListDog(_ *pb.ListDogRequest, stream pb.DogService_ListDogServer) error {
	cur, err := collection.Find(context.Background(), primitive.D{{}})
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error: %v", err),
		)
	}
	defer cur.Close(context.Background()) // should handle err
	for cur.Next(context.Background()) {
		data := &dogData{}
		err := cur.Decode(data)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error while decoding data from MongoDB: %v", err),
			)
		}
		stream.Send(&pb.ListDogResponse{Dog: dataToDogPb(data)}) // should handle err
	}
	if err := cur.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error: %v", err),
		)
	}
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client, err := connectToMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	collection = client.Database("mydb").Collection("blog")

	fmt.Println("Dog API started")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)
	pb.RegisterDogServiceServer(s, &server{})

	reflection.Register(s)

	go func() {
		fmt.Println("Starting server")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	fmt.Println("Stopping the listener")
	lis.Close()
	fmt.Println("Closing mongodb connection")
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatalf("Error on disconnection with MongoDB: %v", err)
	}
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("End of program")
}

func connectToMongoDB() (*mongo.Client, error) {
	var cred options.Credential

	cred.AuthSource = "admin"
	cred.Username = "mongoadmin"
	cred.Password = "password"
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(cred)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Connect(context.TODO())
	return client, err
}

func dataToDogPb(data *dogData) *pb.Dog {
	return &pb.Dog{
		Id:       data.ID.Hex(),
		RegName:  data.RegName,
		CallName: data.CallName,
		Gender:   data.Gender,
	}
}
