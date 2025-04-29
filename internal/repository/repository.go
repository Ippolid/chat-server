package repository

import (
	"context"
	"github.com/Ippolid/chat-server/internal/model"
)

// AuthRepository интерфейс для работы с репозиторием
type AuthRepository interface {
	CreateRequest(ctx context.Context, chat model.Chats) (int64, error)
	DeleteChat(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, message model.MessageInfo) error
	MakeLog(ctx context.Context, log model.Log) error
}
