package chatserver

import (
	"context"
	"fmt"
	"time"

	"github.com/Ippolid/chat-server/internal/model"
)

func (s *serv) Create(ctx context.Context, info *model.Chats) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.chatserverRepository.CreateRequest(ctx, *info)
		if errTx != nil {
			return errTx
		}

		err := s.chatserverRepository.MakeLog(ctx, model.Log{
			Method:    "Create",
			CreatedAt: time.Now(),
			Ctx:       fmt.Sprintf("%v", ctx),
		})

		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}
