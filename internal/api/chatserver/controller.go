package chatserver

import (
	"github.com/Ippolid/chat-server/internal/service"
	"github.com/Ippolid/chat-server/pkg/chatserver_v1"
)

// Controller реализует интерфейс ChatserverV1Server и ChatserverService
type Controller struct {
	chatserver_v1.UnimplementedChatV1Server
	chatserverService service.ChatServerService
}

// NewController создает новый экземпляр Controller
func NewController(chatserverService service.ChatServerService) *Controller {
	return &Controller{
		chatserverService: chatserverService,
	}
}
