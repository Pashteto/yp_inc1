package repos

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// Repository represent the repositories
type SetterGetter interface {
	Set(ctx context.Context, key string, value interface{}, exp time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Ping(ctx context.Context) error
}

// repository represent the repository model
type repository struct {
	Client *redis.Cmdable
}

// NewRedisRepository will create an object that represent the Repository interface
func NewRedisRepository(Client redis.Cmdable) SetterGetter {
	return &repository{&Client}
}

// Set attaches the redis repository and set the data
func (r *repository) Set(ctx context.Context, key string, value interface{}, exp time.Duration) error {

	return (*r.Client).Set(ctx, key, value, exp).Err()
}

// Get attaches the redis repository and get the data
func (r *repository) Get(ctx context.Context, key string) (string, error) {
	get := (*r.Client).Get(ctx, key)
	return get.Result()
}

func (r *repository) Ping(ctx context.Context) error {
	return (*r.Client).Ping(ctx).Err()
}
