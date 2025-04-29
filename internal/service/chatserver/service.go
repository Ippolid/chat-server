package chatserver

import (
	"github.com/Ippolid/chat-server/internal/client/db"
	"github.com/Ippolid/chat-server/internal/repository"
	"github.com/Ippolid/chat-server/internal/service"
)

type serv struct {
	chatserverRepository repository.ChatServerRepository
	txManager            db.TxManager
}

// NewService создает новый экземпляр AuthService
func NewService(
	chatserverRepository repository.ChatServerRepository,
	txManager db.TxManager,
) service.ChatServerService {
	return &serv{
		chatserverRepository: chatserverRepository,
		txManager:            txManager,
	}
}
