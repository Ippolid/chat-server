package access

import (
	"github.com/Ippolid/chat-server/internal/service"
)

// AuthInterceptor структура реализующая интерцептор для авторизации.
type AuthInterceptor struct {
	AccessService service.Access
}

// NewAuthInterceptor конструктор для структуры реализующей интерцептор для авторизации.
func NewAuthInterceptor(accessService service.Access) *AuthInterceptor {
	return &AuthInterceptor{
		AccessService: accessService,
	}
}
