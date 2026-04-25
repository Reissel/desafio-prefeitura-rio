package logic

import (
	"context"
	"database/sql"
	"desafio-prefeitura-rio/model"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type DLQStore interface {
	Save(ctx context.Context, notification model.Notification) (int, error)
}

type PostgresDLQRepo struct {
	DB *sql.DB
}

type RedisDLQRepo struct {
	Client *redis.Client
	Key    string
}

func (p *PostgresDLQRepo) Save(ctx context.Context, n model.Notification) (int, error) {
	jsonNotification, err := json.Marshal(n)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO notifications_dlq (body) VALUES ($1) RETURNING id`

	var createdId int
	err = p.DB.QueryRowContext(ctx, query, jsonNotification).Scan(&createdId)
	return createdId, err
}

func (r *RedisDLQRepo) Save(ctx context.Context, n model.Notification) (int, error) {
	jsonNotification, err := json.Marshal(n)
	if err != nil {
		return 0, err
	}

	err = r.Client.RPush(ctx, r.Key, jsonNotification).Err()
	if err != nil {
		return 0, err
	}

	len, err := r.Client.LLen(ctx, r.Key).Result()
	return int(len), err
}
