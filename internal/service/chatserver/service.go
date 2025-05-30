package chatserver

import (
	"github.com/Ippolid/chat-server/internal/repository"
	"github.com/Ippolid/chat-server/internal/service"
	"github.com/Ippolid/platform_libary/pkg/db"
)

type serv struct {
	chatserverRepository repository.ChatServerRepository
	txManager            db.TxManager
}

// NewService создает новый экземпляр ChatserverService
func NewService(
	chatserverRepository repository.ChatServerRepository,
	txManager db.TxManager,
) service.ChatServerService {
	return &serv{
		chatserverRepository: chatserverRepository,
		txManager:            txManager,
	}
}
