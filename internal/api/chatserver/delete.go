package chatserver

import (
	"context"
	"log"

	"github.com/Ippolid/chat-server/pkg/chatserver_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Delete реализует метод удаления чата
func (i *Controller) Delete(ctx context.Context, req *chatserver_v1.DeleteRequest) (*emptypb.Empty, error) {
	err := i.chatserverService.Delete(ctx, req.GetId())

	if err != nil {
		log.Printf("failed to delete user: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &emptypb.Empty{}, nil
}
