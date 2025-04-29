package chatserver

import (
	"context"
	"fmt"
	"github.com/Ippolid/chat-server/internal/model"
	"time"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.chatserverRepository.DeleteChat(ctx, id)
		if errTx != nil {
			return errTx
		}

		err := s.chatserverRepository.MakeLog(ctx, model.Log{
			Method:    "Delete",
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
