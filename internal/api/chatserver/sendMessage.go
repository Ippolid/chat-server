package chatserver

import (
	"context"
	"log"

	"github.com/Ippolid/chat-server/internal/conventer"
	"github.com/Ippolid/chat-server/pkg/chatserver_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// SendMessage реализует метод получения пользователя по ID
func (i *Controller) SendMessage(ctx context.Context, req *chatserver_v1.SendMessageRequest) (*emptypb.Empty, error) {
	err := i.chatserverService.SendMessage(ctx, conventer.ToMessageFromService(req))
	if err != nil {
		log.Printf("failed to send message: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to send message: %v", err)
	}

	return &emptypb.Empty{}, nil
}
