package query

import (
	"context"
	"fmt"
	"time"
)

// CreateRequest создает новый чат
func (d *Db) CreateRequest(ctx context.Context, users []string) (int, error) {
	var id int
	err := d.db.QueryRow(ctx,
		"INSERT INTO chats (users) VALUES ($1) RETURNING id", users).Scan(&id)

	if err != nil {

		return 0, err
	}

	return id, nil
}

// DeleteChat удаляет чат по id
func (d *Db) DeleteChat(ctx context.Context, id int) (context.Context, error) {
	_, err := d.db.Exec(ctx, "DELETE FROM chats WHERE id=$1", id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete chat: %w", err)
	}
	return nil, nil
}

// SendMessage отправляет сообщение в чат
func (d *Db) SendMessage(ctx context.Context, info MessageInfo, time time.Time) (context.Context, error) {
	_, err := d.db.Exec(ctx, "INSERT INTO messages (from_user, text, sent_at) VALUES ($1, $2, $3)", info.From, info.Text, time)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return nil, nil

}
