package models

import (
	"context"
	"errors"
)

var ErrConfigNotFound = errors.New("config not found")

type Config struct {
	ID                 int  `json:"id"`
	Debug              bool `json:"debug"`
	AttemptNotify      bool `json:"attempt_notify"`
	MaxMessageLength   int  `json:"max_message_length"`
	MaxConcurrentSyncs int  `json:"max_concurrent_syncs"`
}

type ConfigRepository interface {
	Get(ctx context.Context) (Config, error)
	InsertDefault(ctx context.Context) (Config, error)
	Upsert(ctx context.Context, c Config) (Config, error)
}
