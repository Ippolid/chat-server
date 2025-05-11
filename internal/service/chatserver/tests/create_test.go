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

const Ctxstring = "context.Background"

func TestCreate(t *testing.T) {
	type chatserverRepositoryMockFunc func(mc *minimock.Controller) repository.ChatServerRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx  context.Context
		info *model.Chats
	}

	var (
		ctx     = context.Background()
		id      = gofakeit.Int64()
		repoErr = fmt.Errorf("repo error")
		logErr  = fmt.Errorf("log error")
		chats   = &model.Chats{
			Users:     []string{gofakeit.Username(), gofakeit.Username()},
			CreatedAt: time.Now(),
		}
	)

	tests := []struct {
		name                     string
		args                     args
		wantID                   int64
		wantErr                  error
		chatserverRepositoryMock chatserverRepositoryMockFunc
		txManagerMock            txManagerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:  ctx,
				info: chats,
			},
			wantID:  id,
			wantErr: nil,
			chatserverRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
				mock := repoMocks.NewChatServerRepositoryMock(mc)
				mock.CreateRequestMock.Expect(ctx, *chats).Return(id, nil)
				mock.MakeLogMock.Set(func(_ context.Context, log model.Log) error {
					if log.Method != "Create" || log.Ctx != Ctxstring {
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
			name: "CreateRequest error case",
			args: args{
				ctx:  ctx,
				info: chats,
			},
			wantID:  0,
			wantErr: repoErr,
			chatserverRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
				mock := repoMocks.NewChatServerRepositoryMock(mc)
				mock.CreateRequestMock.Expect(ctx, *chats).Return(int64(0), repoErr)
				// MakeLog не должен вызываться при ошибке CreateRequest
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := servMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) error {
					return f(ctx) // Ошибка из CreateRequest должна быть возвращена здесь
				})
				return mock
			},
		},
		{
			name: "MakeLog error case",
			args: args{
				ctx:  ctx,
				info: chats,
			},
			wantID:  0,
			wantErr: logErr,
			chatserverRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
				mock := repoMocks.NewChatServerRepositoryMock(mc)
				mock.CreateRequestMock.Expect(ctx, *chats).Return(id, nil)
				mock.MakeLogMock.Set(func(_ context.Context, log model.Log) error {
					if log.Method != "Create" || log.Ctx != Ctxstring {
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
			t.Parallel() // Можно добавить для параллельного выполнения тестов
			mc := minimock.NewController(t)
			defer mc.Finish() // Убедимся, что все ожидания были выполнены

			chatRepoMock := tt.chatserverRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)
			service := chatserver.NewService(chatRepoMock, txManagerMock)

			gotID, err := service.Create(tt.args.ctx, tt.args.info)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.wantID, gotID)
		})
	}
}
