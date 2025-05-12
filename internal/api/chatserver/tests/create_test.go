package chatserver_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/Ippolid/chat-server/internal/api/chatserver"
	"github.com/Ippolid/chat-server/internal/model"
	"github.com/Ippolid/chat-server/internal/service"
	"github.com/Ippolid/chat-server/internal/service/mocks"
	desc "github.com/Ippolid/chat-server/pkg/chatserver_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestController_Create(t *testing.T) {
	type chatserverServiceMockFunc func(mc *minimock.Controller) service.ChatServerService // Используем интерфейс сервиса чата

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()

		id = gofakeit.Int64()

		serviceErr = fmt.Errorf("service error")

		usernames = []string{gofakeit.Username(), gofakeit.Username()}
		req       = &desc.CreateRequest{
			Usernames: usernames,
		}
	)

	tests := []struct {
		name                  string
		args                  args
		wantResp              *desc.CreateResponse
		wantErr               error
		chatserverServiceMock chatserverServiceMockFunc
	}{
		{
			name: "success",
			args: args{ctx: ctx, req: req},
			wantResp: &desc.CreateResponse{
				Id: id,
			},
			wantErr: nil,
			chatserverServiceMock: func(mc *minimock.Controller) service.ChatServerService {
				mock := mocks.NewChatServerServiceMock(mc)
				mock.CreateMock.Set(func(c context.Context, info *model.Chats) (int64, error) {
					require.Equal(t, ctx, c)
					require.ElementsMatch(t, usernames, info.Users)
					return id, nil
				})
				return mock
			},
		},
		{
			name:     "service error",
			args:     args{ctx: ctx, req: req},
			wantResp: nil,
			wantErr:  serviceErr,
			chatserverServiceMock: func(mc *minimock.Controller) service.ChatServerService {
				mock := mocks.NewChatServerServiceMock(mc) // Мок сервиса чата
				mock.CreateMock.Set(func(c context.Context, info *model.Chats) (int64, error) {
					// Используем 't' из замыкания t.Run
					require.Equal(t, ctx, c)
					require.ElementsMatch(t, usernames, info.Users) // Сравниваем только Usernames
					return int64(0), serviceErr
				})
				return mock
			},
		},
		{
			name:     "service error on nil model from converter",              // Случай, если конвертер вернул nil
			args:     args{ctx: ctx, req: &desc.CreateRequest{Usernames: nil}}, // Пример запроса, который может дать nil
			wantResp: nil,
			wantErr:  serviceErr, // Ожидаем ошибку от сервиса
			chatserverServiceMock: func(mc *minimock.Controller) service.ChatServerService {
				mock := mocks.NewChatServerServiceMock(mc)
				mock.CreateMock.Set(func(c context.Context, info *model.Chats) (int64, error) {
					// Используем 't' из замыкания t.Run
					require.Equal(t, ctx, c)
					require.Nil(t, info.Users)
					return int64(0), serviceErr
				})
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { // 't' здесь - это *testing.T для каждого подтеста
			// mc - это контроллер моков от minimock.
			// minimock.NewController(t) инициализирует контроллер с текущим тестовым контекстом (*testing.T).
			mc := minimock.NewController(t)
			defer mc.Finish()

			chatServiceMock := tt.chatserverServiceMock(mc)
			api := chatserver.NewController(chatServiceMock) // Используем контроллер чата
			resp, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.wantResp, resp)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
