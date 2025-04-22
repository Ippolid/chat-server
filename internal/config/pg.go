package config

import (
	"errors"
	"os"
)

const (
	// DSNKey - ключ для строки подключения к PostgreSQL
	DSNKey = "PG_DSN"
)

// PGConfig представляет интерфейс для получения строки подключения к PostgreSQL
type PGConfig interface {
	DSN() string
}

type pgConfig struct {
	dsn string
}

// NewPGConfig создает новый экземпляр PGConfig, получая DSN из переменной окружения
func NewPGConfig() (PGConfig, error) {
	dsn := os.Getenv(DSNKey)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}

// возвратит DSN для подключения к PostgreSQL
func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
