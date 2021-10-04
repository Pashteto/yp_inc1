package repos

import (
	"context"
	"log"
	"sort"
	"time"

	"github.com/go-redis/redis/v8"
)

// Repository represent the repositories
type SetterGetter interface {
	Ping(ctx context.Context) error
	FlushAllKeys(ctx context.Context) error

	SetValueByKeyAndUser(ctx context.Context, key, UserHash, URL string, exp time.Duration) error
	GetValueByKey(ctx context.Context, key, UserHash string, UserList *[]string) (string, error)

	ListAllKeysByUser(ctx context.Context, UserHash string) (map[string]string, error)
}

// repository represent the repository model
type repository struct {
	Client redis.Cmdable
}

// NewRedisRepository will create an object that represent the Repository interface
func NewRedisRepository(Client redis.Cmdable) SetterGetter {
	return &repository{Client}
}

func (r *repository) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

func (r *repository) FlushAllKeys(ctx context.Context) error {
	return r.Client.FlushAll(ctx).Err()
}

func (r *repository) SetValueByKeyAndUser(ctx context.Context, key, User, URL string, exp time.Duration) error {
	err := r.Client.HSet(ctx, User, key, URL).Err()
	//	data, _ := r.Client.HGetAll(ctx, UserHash).Result()
	//	log.Println(UserHash, "\tin HSET\t", data)
	what := r.Client.Expire(ctx, User, exp).Err()
	if what != nil {
		return what
	}
	return err
}

// Get stored URL value by Key w/o username
func (r *repository) GetValueByKey(ctx context.Context, key, User string, UserList *[]string) (string, error) {
	//	data, _ := r.Client.HGetAll(ctx, UserHash).Result()
	//	log.Println(UserHash, "\tin r.Client.HGetAll\t", data)
	var res string
	for _, UserFromList := range *UserList {
		res, _ = r.Client.HGet(ctx, UserFromList, key).Result()
		if res != "" {
			log.Println(User, "\tr.Client.HGet(ctx, UserFromList, key).Result()\t", res)
			return res, nil
		}
	}
	return "", nil
}

func (r *repository) ListAllKeysByUser(ctx context.Context, User string) (map[string]string, error) {
	return r.Client.HGetAll(ctx, User).Result()
}

func Contains(s *[]string, searchterm string) bool {
	i := sort.SearchStrings(*s, searchterm)
	return i < len(*s) && (*s)[i] == searchterm
}
