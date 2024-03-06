package redis

import (
	"context"
	"fmt"
	"git.spbec-mining.ru/arxon31/sambaMW/pkg/logger/sl"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"time"
)

type Redis struct {
	client *redis.Client
	logger *slog.Logger
}

func New(ctx context.Context, logger *slog.Logger, host, port, password string, db int) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	context, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	status := client.Ping(context)
	if status.Err() != nil {
		return &Redis{}, status.Err()
	}

	logger.Debug("connected to redis", slog.String("host", host), slog.String("port", port))

	return &Redis{client: client, logger: logger}, nil
}

func (r *Redis) SaveDirs(ctx context.Context, dirs []string) error {
	for _, dir := range dirs {
		err := r.client.Set(ctx, dir, "empty", 0).Err()
		if err != nil {
			r.logger.Error("failed to save dir", sl.Err(err))
			return err
		}
	}
	return nil
}

func (r *Redis) GetDirs(ctx context.Context) ([]string, error) {
	dirs, err := r.client.Keys(ctx, "*").Result()
	if err != nil {
		r.logger.Error("failed to get dirs", sl.Err(err))
		return nil, err
	}
	return dirs, nil
}

func (r *Redis) DeleteEmptyDir(ctx context.Context, dir string) error {

	err := r.client.Del(ctx, dir).Err()
	if err != nil {
		r.logger.Error("failed to delete dir", sl.Err(err))
		return err
	}

	return nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}
