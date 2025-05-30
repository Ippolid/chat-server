package conventer

import (
	"github.com/Ippolid/chat-server/internal/model"
	"github.com/Ippolid/chat-server/pkg/chatserver_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToChatFromService преобразует запрос на создание чата в модель чата
func ToChatFromService(req *chatserver_v1.CreateRequest) *model.Chats {
	users := req.GetUsernames()
	createdAt := timestamppb.Now()
	user := model.Chats{
		Users:     users,
		CreatedAt: createdAt.AsTime(),
	}

	return &user
}

// ToMessageFromService преобразует запрос на отправку сообщения в модель сообщения
func ToMessageFromService(req *chatserver_v1.SendMessageRequest) *model.MessageInfo {
	messageInfo := model.MessageInfo{
		From:   req.GetMessage().From,
		Text:   req.GetMessage().Text,
		SentAt: req.GetTimestamp().AsTime(),
	}

	return &messageInfo
}
