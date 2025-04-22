package server

import (
	"context"
	"fmt"
	"log"

	"github.com/Ippolid/chat-server/pkg/chatserver_v1"
	"github.com/Ippolid/chat-server/postgres/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Create chat
func (s *Server) Create(_ context.Context, req *chatserver_v1.CreateRequest) (*chatserver_v1.CreateResponse, error) {

	ctx := context.Background()
	id, err := s.db.CreateRequest(ctx, req.GetUsernames())
	if err != nil {
		log.Printf("failed to get user: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get users: %v", err)
	}
	fmt.Printf("Users: %+v", id)

	return &chatserver_v1.CreateResponse{
		Id: int64(id),
	}, nil
}

// Delete chat
func (s *Server) Delete(_ context.Context, req *chatserver_v1.DeleteRequest) (*emptypb.Empty, error) {
	//чето делается
	ctx := context.Background()
	_, err := s.db.DeleteChat(ctx, int(req.GetId()))
	if err != nil {
		log.Printf("failed to delete chat: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete chat: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// SendMessage отправляет сообщение в чат
func (s *Server) SendMessage(_ context.Context, req *chatserver_v1.SendMessageRequest) (*emptypb.Empty, error) {
	//чето делается
	from := req.GetMessage().From
	text := req.GetMessage().Text

	message := query.MessageInfo{
		From: from,
		Text: text,
	}
	ctx := context.Background()
	_, err := s.db.SendMessage(ctx, message, req.Timestamp.AsTime())
	if err != nil {
		log.Printf("failed to send message: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to send message: %v", err)
	}

	return &emptypb.Empty{}, nil
}
