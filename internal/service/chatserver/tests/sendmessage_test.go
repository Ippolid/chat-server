package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Ippolid/chat-server/internal/model"
	"github.com/Ippolid/chat-server/internal/repository"
	repoMocks "github.com/Ippolid/chat-server/internal/repository/mocks"
	"github.com/Ippolid/chat-server/internal/service/chatserver"
	servMocks "github.com/Ippolid/chat-server/internal/service/mocks"
	"github.com/Ippolid/platform_libary/pkg/db"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestSendMessage(t *testing.T) {
	type chatserverRepositoryMockFunc func(mc *minimock.Controller) repository.ChatServerRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx     context.Context
		message *model.MessageInfo
	}

	var (
		ctx     = context.Background()
		repoErr = fmt.Errorf("repo error")
		logErr  = fmt.Errorf("log error")
		message = &model.MessageInfo{
			From:   gofakeit.Username(),
			Text:   gofakeit.Sentence(10),
			SentAt: time.Now(),
		}
	)

	tests := []struct {
		name                     string
		args                     args
		wantErr                  error
		chatserverRepositoryMock chatserverRepositoryMockFunc
		txManagerMock            txManagerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:     ctx,
				message: message,
			},
			wantErr: nil,
			chatserverRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
				mock := repoMocks.NewChatServerRepositoryMock(mc)
				mock.SendMessageMock.Expect(ctx, *message).Return(nil)
				mock.MakeLogMock.Set(func(_ context.Context, log model.Log) error {
					if log.Method != "SendMessage" || log.Ctx != Ctxstring {
						return fmt.Errorf("unexpected log entry: %+v", log)
					}
					return nil
				})
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := servMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) error {
					return f(ctx)
				})
				return mock
			},
		},
		{
			name: "SendMessage repo error case",
			args: args{
				ctx:     ctx,
				message: message,
			},
			wantErr: repoErr,
			chatserverRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
				mock := repoMocks.NewChatServerRepositoryMock(mc)
				mock.SendMessageMock.Expect(ctx, *message).Return(repoErr)
				// MakeLog не должен вызываться при ошибке SendMessage
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := servMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) error {
					return f(ctx) // Ошибка из SendMessage должна быть возвращена здесь
				})
				return mock
			},
		},
		{
			name: "MakeLog error case",
			args: args{
				ctx:     ctx,
				message: message,
			},
			wantErr: logErr,
			chatserverRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
				mock := repoMocks.NewChatServerRepositoryMock(mc)
				mock.SendMessageMock.Expect(ctx, *message).Return(nil)
				mock.MakeLogMock.Set(func(_ context.Context, log model.Log) error {
					if log.Method != "SendMessage" || log.Ctx != Ctxstring {
						return fmt.Errorf("unexpected log entry: %+v", log)
					}
					return logErr
				})
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := servMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) error {
					return f(ctx) // Ошибка из MakeLog должна быть возвращена здесь
				})
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)
			defer mc.Finish()

			chatRepoMock := tt.chatserverRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)
			service := chatserver.NewService(chatRepoMock, txManagerMock)

			err := service.SendMessage(tt.args.ctx, tt.args.message)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
