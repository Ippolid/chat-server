package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Ippolid/chat-server/grpc/pkg/chatserver_v1"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50051

type server struct {
	chatserver_v1.UnimplementedAuthV1Server
}

// Get ...

func (s *server) Create(_ context.Context, req *chatserver_v1.CreateRequest) (*chatserver_v1.CreateResponse, error) {
	//чето делается
	fmt.Printf("name +%v\n", req.Usernames)

	return &chatserver_v1.CreateResponse{
		Id: gofakeit.Int64(),
	}, nil
}

func (s *server) Delete(_ context.Context, req *chatserver_v1.DeleteRequest) (*emptypb.Empty, error) {
	//чето делается
	fmt.Printf("User id: %d", req.GetId())

	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(_ context.Context, req *chatserver_v1.SendMessageRequest) (*emptypb.Empty, error) {
	//чето делается
	fmt.Printf("User id: %d", req.Message, req.Timestamp)

	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	chatserver_v1.RegisterAuthV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
