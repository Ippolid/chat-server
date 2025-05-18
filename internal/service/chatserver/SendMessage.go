package chatserver

import (
	"context"
	"fmt"
	"time"

	"github.com/Ippolid/chat-server/internal/model"
)

func (s *serv) SendMessage(ctx context.Context, info *model.MessageInfo) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.chatserverRepository.SendMessage(ctx, *info)
		if errTx != nil {
			return errTx
		}

		err := s.chatserverRepository.MakeLog(ctx, model.Log{
			Method:    "SendMessage",
			CreatedAt: time.Now(),
			Ctx:       fmt.Sprintf("%v", ctx),
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
