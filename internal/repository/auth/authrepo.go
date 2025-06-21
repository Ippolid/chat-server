package auth

import (
	"github.com/Ippolid/auth/pkg/auth_v1"
)

// RepoAccess структура реализующая репозиторий для доступа к сервису авторизации.
type RepoAccess struct {
	client auth_v1.AuthClient
}

// NewAccessRepo конструктор для структуры реализующей репозиторий для доступа к сервису авторизации.
func NewAccessRepo(client auth_v1.AuthClient) *RepoAccess {
	return &RepoAccess{client: client}
}
