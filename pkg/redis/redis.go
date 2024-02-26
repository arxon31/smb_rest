package redis

import (
	"context"
	"fmt"
	"git.spbec-mining.ru/arxon31/sambaMW/pkg/logger/sl"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type Redis struct {
	client *redis.Client
	logger *slog.Logger
}

func New(ctx context.Context, host, port string, logger *slog.Logger) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "",
		DB:       0,
	})

	status := client.Ping(ctx)
	if status.Err() != nil {
		return &Redis{}
	}

	return &Redis{client: client, logger: logger}
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

func (r *Redis) DeleteDir(ctx context.Context, dir string) error {

	err := r.client.Del(ctx, dir).Err()
	if err != nil {
		r.logger.Error("failed to delete dir", sl.Err(err))
		return err
	}

	return nil
}
