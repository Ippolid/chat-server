package chatserver_tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Ippolid/chat-server/internal/api/chatserver"
	"github.com/Ippolid/chat-server/internal/model"
	"github.com/Ippolid/chat-server/internal/service"
	"github.com/Ippolid/chat-server/internal/service/mocks"
	desc "github.com/Ippolid/chat-server/pkg/chatserver_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestController_SendMessage(t *testing.T) {
	type chatserverServiceMockFunc func(mc *minimock.Controller) service.ChatServerService

	type args struct {
		ctx context.Context
		req *desc.SendMessageRequest
	}

	var (
		ctx = context.Background()

		fromUsername = gofakeit.Username()
		text         = gofakeit.Sentence(10)
		timestamp    = timestamppb.New(time.Now())

		serviceErr = fmt.Errorf("service error")

		req = &desc.SendMessageRequest{
			Message: &desc.MessageInfo{
				From: fromUsername,
				Text: text,
			},
			Timestamp: timestamp,
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
				mock.SendMessageMock.Set(func(c context.Context, msg *model.MessageInfo) error {
					require.Equal(t, ctx, c)
					require.Equal(t, fromUsername, msg.From)
					require.Equal(t, text, msg.Text)
					return nil
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
				mock := mocks.NewChatServerServiceMock(mc)
				mock.SendMessageMock.Set(func(c context.Context, msg *model.MessageInfo) error {
					require.Equal(t, ctx, c)
					require.Equal(t, fromUsername, msg.From)
					require.Equal(t, text, msg.Text)
					return serviceErr
				})
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

			resp, err := api.SendMessage(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.wantResp, resp)
			if tt.wantErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr.Error())
			}
		})
	}
}
