package service

import (
	"context"

	"github.com/Ippolid/chat-server/internal/model"
)

// ChatServerService определяет интерфейс для работы с чатом
type ChatServerService interface {
	Create(ctx context.Context, info *model.Chats) (int64, error)
	Delete(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, info *model.MessageInfo) error
}

type Access interface {
	Access(ctx context.Context, path string) error
}
