package auth

import (
	"github.com/Ippolid/auth/pkg/auth_v1"
)

type RepoAccess struct {
	client auth_v1.AuthClient
}

func NewAccessRepo(client auth_v1.AuthClient) *RepoAccess {
	return &RepoAccess{client: client}
}
