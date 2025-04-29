package chat_server

import (
	"context"
	"fmt"

	"github.com/Ippolid/chat-server/internal/client/db"
	"github.com/Ippolid/chat-server/internal/model"
	"github.com/Ippolid/chat-server/internal/repository"
	sq "github.com/Masterminds/squirrel"
)

const (
	tableChatName    = "chats"
	tableMessageName = "messages"
	tableLogName     = "logs"

	idColumn        = "id"
	usersColumn     = "users"
	createdAtColumn = "created_at"

	fromColumn   = "from_user"
	textColumn   = "text"
	sentAtColumn = "sent_at"

	methodColumn = "method_name"
	ctxColumn    = "ctx"
)

type repo struct {
	db db.Client
}

// NewRepository создает новый экземпляр репозитория
func NewRepository(db db.Client) repository.ChatServerRepository {
	return &repo{db: db}
}

// CreateRequest создает новый чат
func (r *repo) CreateRequest(ctx context.Context, chat model.Chats) (int64, error) {
	builder := sq.Insert(tableChatName).
		PlaceholderFormat(sq.Dollar).
		Columns(usersColumn).
		Values(chat.Users).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "chatserver_repository.CreateRequest",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// DeleteChat удаляет чат по ID из базы данных
func (r *repo) DeleteChat(ctx context.Context, id int64) error {
	builder := sq.Delete(tableChatName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chatserver_repository.Delete_chat",
		QueryRaw: query,
	}

	tag, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	// Опционально: проверка, что запись действительно была удалена
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("chat with id %d not found", id)
	}

	return nil
}

// DeleteUser удаляет пользователя по ID из базы данных
func (r *repo) SendMessage(ctx context.Context, message model.MessageInfo) error {
	builder := sq.Insert(tableMessageName).
		PlaceholderFormat(sq.Dollar).
		Columns(fromColumn, textColumn, sentAtColumn).
		Values(message.From, message.Text, message.SentAt)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chatserver_repository.Send_message",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) MakeLog(ctx context.Context, info model.Log) error {
	builder := sq.Insert(tableLogName).
		PlaceholderFormat(sq.Dollar).
		Columns(methodColumn, createdAtColumn, ctxColumn).
		Values(info.Method, info.CreatedAt, info.Ctx)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chatserver_repository.Log",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}
