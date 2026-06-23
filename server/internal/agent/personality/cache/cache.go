package cache

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
)

const PersonalityKey = "agent:personality:default"

type Store interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte) error
	Delete(ctx context.Context, key string) error
}

type Redis struct {
	client *redis.Client
}

func NewRedis(url string) (*Redis, error) {
	opts, err := redis.ParseURL(strings.TrimSpace(url))
	if err != nil {
		return nil, err
	}
	return &Redis{client: redis.NewClient(opts)}, nil
}

func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

func (r *Redis) Set(ctx context.Context, key string, value []byte) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

type Noop struct{}

func (Noop) Get(context.Context, string) ([]byte, error) { return nil, nil }
func (Noop) Set(context.Context, string, []byte) error  { return nil }
func (Noop) Delete(context.Context, string) error        { return nil }
