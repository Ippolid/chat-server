package access

import (
	"github.com/Ippolid/chat-server/internal/repository"
	"github.com/Ippolid/chat-server/internal/service"
)

// ServiceAccess структура реализующая сервис интерцептора.
type ServiceAccess struct {
	accessRepository service.Access
}

// NewAccessService конструктор для сервиса интерцептора.
func NewAccessService(accessRepository repository.Access) service.Access {
	return &ServiceAccess{
		accessRepository: accessRepository,
	}
}
