package query

import (
	"github.com/jackc/pgx/v5"
)

// Db представляет структуру базы данных
type Db struct {
	db *pgx.Conn
}

// NewDb создает новый экземпляр Db с подключением к базе данных
func NewDb(db *pgx.Conn) *Db {
	return &Db{
		db: db,
	}
}

// MessageInfo представляет структуру информации о сообщении
type MessageInfo struct {
	From string
	Text string
}
