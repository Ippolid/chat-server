package server

import (
	"github.com/Ippolid/chat-server/pkg/chatserver_v1"
	"github.com/Ippolid/chat-server/postgres/query"
	"github.com/jackc/pgx/v5"
)

// Server представляет структуру сервера
type Server struct {
	chatserver_v1.UnimplementedChatV1Server
	db *query.Db
}

// NewServer создает новый экземпляр сервера с подключением к базе данных
func NewServer(db *pgx.Conn) *Server {
	return &Server{
		db: query.NewDb(db),
	}
}
