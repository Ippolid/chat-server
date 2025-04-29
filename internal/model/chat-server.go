package model

import (
	"time"
)

// MessageInfo Message структура сообщения
type MessageInfo struct {
	From   string    `db:"from-user"`
	Text   string    `db:"text"`
	SentAt time.Time `db:"sent_at"`
}

// Chats Chat структура чата
type Chats struct {
	Users     []string  `db:"users"`
	CreatedAt time.Time `db:"created_at"`
}

// Log структура для хранения логов
type Log struct {
	Method    string    `db:"method_name"`
	CreatedAt time.Time `db:"created_at"`
	Ctx       string    `db:"ctx"`
}
