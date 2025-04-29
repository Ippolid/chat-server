package chatserver

import (
	"context"
	"github.com/Ippolid/chat-server/internal/conventer"
	"github.com/Ippolid/chat-server/pkg/chatserver_v1"

	"log"
)

// Create реализует метод создания чата
func (i *Controller) Create(ctx context.Context, req *chatserver_v1.CreateRequest) (*chatserver_v1.CreateResponse, error) {
	id, err := i.chatserverService.Create(ctx, conventer.ToChatFromService(req))
	if err != nil {
		return nil, err
	}

	log.Printf("inserted note with id: %d", id)

	return &chatserver_v1.CreateResponse{
		Id: id,
	}, nil
}
