package repos

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// Repository represent the repositories
type SetterGetter interface {
	Set(ctx context.Context, key string, value UserAndString, exp time.Duration) error
	Get(ctx context.Context, key string) (UserAndString, error)
	Ping(ctx context.Context) error
	ListAllKeys(ctx context.Context) ([]string, error)
	FlushAllKeys(ctx context.Context) error

	SetHash(ctx context.Context, key, UserHash, URL string, exp time.Duration) error
	GetHash(ctx context.Context, key, UserHash string) (string, error)
	ListAllKeysHashed(ctx context.Context, UserHash string) (map[string]string, error)
}

// repository represent the repository model
type repository struct {
	Client redis.Cmdable
}

// NewRedisRepository will create an object that represent the Repository interface
func NewRedisRepository(Client redis.Cmdable) SetterGetter {
	return &repository{Client}
}

// Set attaches the redis repository and set the data
func (r *repository) Set(ctx context.Context, key string, value UserAndString, exp time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.Client.Set(ctx, key, p, exp).Err()
}

// Get attaches the redis repository and get the data
func (r *repository) Get(ctx context.Context, key string) (UserAndString, error) {
	var res UserAndString
	p, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return UserAndString{"", ""}, err
	}
	err = json.Unmarshal([]byte(p), &res)
	return res, err
}

func (r *repository) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

func (r *repository) ListAllKeys(ctx context.Context) ([]string, error) {
	return r.Client.Keys(ctx, "*").Result()
}

func (r *repository) FlushAllKeys(ctx context.Context) error {
	return r.Client.FlushAll(ctx).Err()
}

type UserAndString struct {
	User string
	URL  string
}

func (r *repository) SetHash(ctx context.Context, key, UserHash, URL string, exp time.Duration) error {
	err := r.Client.HSet(ctx, UserHash, key, URL).Err()
	//	data, _ := r.Client.HGetAll(ctx, UserHash).Result()
	//	log.Println(UserHash, "\tin HSET\t", data)
	what := r.Client.Expire(ctx, UserHash, exp).Err()
	if what != nil {
		return what
	}
	return err
}

// Get attaches the redis repository and get the data
func (r *repository) GetHash(ctx context.Context, key, UserHash string) (string, error) {
	data, _ := r.Client.HGetAll(ctx, UserHash).Result()
	log.Println(UserHash, "\tin r.Client.HGetAll\t", data)

	data1, _ := r.Client.Keys(ctx, "*").Result()
	log.Println(UserHash, "\tin r.Client.Keys\t", data1)
	var res string
	for _, dataa := range data1 {
		res, _ = r.Client.HGet(ctx, dataa, key).Result()
		if res != "" {
			log.Println(UserHash, "\tr.Client.HGet(ctx, dataa, key).Result()\t", res)
			return res, nil
		}
	}
	return "", nil
}

func (r *repository) ListAllKeysHashed(ctx context.Context, UserHash string) (map[string]string, error) {
	return r.Client.HGetAll(ctx, UserHash).Result()
}
