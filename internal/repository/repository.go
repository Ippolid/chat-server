package repository

import (
	"context"

	"github.com/Ippolid/chat-server/internal/model"
)

// ChatServerRepository интерфейс для слоя репо
type ChatServerRepository interface {
	CreateRequest(ctx context.Context, chat model.Chats) (int64, error)
	DeleteChat(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, message model.MessageInfo) error
	MakeLog(ctx context.Context, log model.Log) error
}

type Access interface {
	Access(ctx context.Context, path string) error
}
