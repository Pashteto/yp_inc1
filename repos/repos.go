package repos

import (
	"context"
	"log"
	"sort"
	"time"

	//	"github.com/go-redis/redis/v8" REDIS
	"github.com/jackc/pgx/v4/pgxpool"
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
	//REDIS
	//Client   redis.Cmdable
	connPool *pgxpool.Pool
}

// NewRedisRepository will create an object that represent the Repository interface
func NewRepoWitnTable(ctx context.Context, connPool *pgxpool.Pool) (SetterGetter, error) {
	sqlCreate := `CREATE TABLE IF NOT EXISTS shorturls(id serial, userid text, keyurl text, longurl text);`
	oi, err := connPool.Exec(ctx, sqlCreate)
	log.Println(oi, err, connPool.Ping(ctx))

	return &repository{connPool}, err
	//{Client, connPool}, err
}

func (r *repository) Ping(ctx context.Context) error {
	return r.connPool.Ping(ctx)
	//REDIS
	//	return r.Client.Ping(ctx).Err()
}

func (r *repository) FlushAllKeys(ctx context.Context) error {
	//REDIS
	// 	r.Client.FlushAll(ctx)
	_, err := r.connPool.Exec(ctx, "TRUNCATE TABLE shorturls;")
	return err
	//REDIS
	//	return r.Client.FlushAll(ctx).Err()
}

func (r *repository) SetValueByKeyAndUser(ctx context.Context, key, User, URL string, exp time.Duration) error {

	//  REDIS
	/*	err1 := r.Client.HSet(ctx, User, key, URL).Err()
		err2 := r.Client.Expire(ctx, User, exp).Err()
		if err1 != nil {
			log.Println(err1, err2)
		}*/
	sqlInsert := "INSERT INTO shorturls (userid , keyurl , longurl) VALUES ($1, $2, $3)"
	_, err := r.connPool.Exec(ctx, sqlInsert, User, key, URL)
	return err
}

// Get stored URL value by Key w/o username
func (r *repository) GetValueByKey(ctx context.Context, key, User string, UserList *[]string) (string, error) {
	queryrow := `SELECT longurl from shorturls WHERE keyurl = $1`
	row := r.connPool.QueryRow(context.Background(), queryrow, key)
	var res string
	//log.Println(row.Scan(&res))
	return res, row.Scan(&res)
	/* REDIS
	var res string
	for _, UserFromList := range *UserList {
		res, _ = r.Client.HGet(ctx, UserFromList, key).Result()
		if res != "" {
			log.Println(User, "\tr.Client.HGet(ctx, UserFromList, key).Result()\t", res)
			return res, nil
		}
	}
	return "", nil
	*/
}

func (r *repository) ListAllKeysByUser(ctx context.Context, User string) (map[string]string, error) {
	AllKeys := make(map[string]string)

	queryrows := `SELECT keyurl , longurl from shorturls WHERE userid = $1`
	rows, err := r.connPool.Query(context.Background(), queryrows, User)
	if err != nil {
		log.Println(err)
	}
	// обязательно закрываем после завершения функции
	defer rows.Close()
	// пробегаем по всем записям
	for rows.Next() {
		var key, URL string
		err = rows.Scan(&key, &URL)
		if err != nil {
			return AllKeys, err
		}
		AllKeys[key] = URL
	}
	log.Println(AllKeys)
	return AllKeys, nil
	//REDIS
	//return r.Client.HGetAll(ctx, User).Result()
}

func Contains(s *[]string, searchterm string) bool {
	i := sort.SearchStrings(*s, searchterm)
	return i < len(*s) && (*s)[i] == searchterm
}
