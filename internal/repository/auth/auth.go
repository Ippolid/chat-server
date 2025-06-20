package auth

import (
	"context"
	"fmt"

	"github.com/Ippolid/auth/pkg/auth_v1"
)

// Access вызываем сервис авторизации для проверки доступа.
func (repo RepoAccess) Access(ctx context.Context, path string) error {
	_, err := repo.client.Check(ctx, &auth_v1.CheckRequest{EndpointAddress: path})
	fmt.Println(err)
	return err
}
