package chatserver_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/Ippolid/chat-server/internal/api/chatserver"
	"github.com/Ippolid/chat-server/internal/service"
	"github.com/Ippolid/chat-server/internal/service/mocks"
	desc "github.com/Ippolid/chat-server/pkg/chatserver_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestController_Delete(t *testing.T) {
	type chatserverServiceMockFunc func(mc *minimock.Controller) service.ChatServerService

	type args struct {
		ctx context.Context
		req *desc.DeleteRequest
	}

	var (
		ctx = context.Background()
		id  = gofakeit.Int64()

		serviceErr = fmt.Errorf("service error")

		req = &desc.DeleteRequest{
			Id: id,
		}
	)

	tests := []struct {
		name                  string
		args                  args
		wantResp              *emptypb.Empty
		wantErr               error
		chatserverServiceMock chatserverServiceMockFunc
	}{
		{
			name:     "success",
			args:     args{ctx: ctx, req: req},
			wantResp: &emptypb.Empty{},
			wantErr:  nil,
			chatserverServiceMock: func(mc *minimock.Controller) service.ChatServerService {
				mock := mocks.NewChatServerServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
				return mock
			},
		},
		{
			name:     "service error",
			args:     args{ctx: ctx, req: req},
			wantResp: nil,
			wantErr:  serviceErr,
			chatserverServiceMock: func(mc *minimock.Controller) service.ChatServerService {
				mock := mocks.NewChatServerServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			chatServiceMock := tt.chatserverServiceMock(mc)
			api := chatserver.NewController(chatServiceMock)

			resp, err := api.Delete(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.wantResp, resp)
			if tt.wantErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err) // Убедимся, что ошибка вообще есть
				require.Contains(t, err.Error(), tt.wantErr.Error())
			}
		})
	}
}
