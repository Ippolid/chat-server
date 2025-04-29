package service

import (
	"context"
	"github.com/Ippolid/chat-server/internal/model"
)

type ChatServerService interface {
	Create(ctx context.Context, info *model.Chats) (int64, error)
	Delete(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, info *model.MessageInfo)
}
