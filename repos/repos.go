package repos

import (
	"context"
	"errors"
	"log"
	"sort"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Repository represent the repositories
type SetterGetter interface {
	Ping(ctx context.Context) error
	FlushAllKeys(ctx context.Context) error

	SetValueByKeyAndUser(ctx context.Context, key, UserHash, URL string, exp time.Duration) error
	GetValueByKey(ctx context.Context, key, UserHash string, UserList *[]string) (string, error)
	GetIDByLong(ctx context.Context, longURL, User string, UserList *[]string) (string, error)

	ListAllKeysByUser(ctx context.Context, UserHash string) (map[string]string, error)
	SetBatch(ctx context.Context, SetsForDB BatchSetsForDB) error
}

// repository represent the repository model
type repository struct {
	connPool *pgxpool.Pool
}

type BatchSetsForDB struct {
	UserID string
	Pairs  []IDShortURL
}
type IDShortURL struct {
	ShortURL string
	LongURL  string
}

// NewRedisRepository will create an object that represent the Repository interface
func NewRepoWitnTable(ctx context.Context, connPool *pgxpool.Pool) (SetterGetter, error) {
	// connPool.Exec(ctx, "DROP TABLE IF EXISTS shorturls")
	sqlCreate := `CREATE TABLE shorturls(id serial, userid text, keyurl text, longurl text, UNIQUE(longurl));`
	_, err := connPool.Exec(ctx, sqlCreate)
	return &repository{connPool}, err
}

func (r *repository) Ping(ctx context.Context) error {
	return r.connPool.Ping(ctx)
}

func (r *repository) FlushAllKeys(ctx context.Context) error {

	_, err := r.connPool.Exec(ctx, "TRUNCATE TABLE shorturls;")
	return err
}

func (r *repository) SetValueByKeyAndUser(ctx context.Context, key, User, URL string, exp time.Duration) error {
	sqlInsert := "INSERT INTO shorturls (userid , keyurl , longurl) VALUES ($1, $2, $3) ON CONFLICT (longurl) DO NOTHING;"
	conflictCheck, err := r.connPool.Exec(ctx, sqlInsert, User, key, URL)
	/*	if err.Error() == pgerrcode.UniqueViolation {
		return errors.New("no way to implement this")
	}*/

	if conflictCheck.String() == "INSERT 0 0" {
		return errors.New(`longURL exists`)
	}
	return err
}

// Get stored URL value by Key w/o username
func (r *repository) GetValueByKey(ctx context.Context, key, User string, UserList *[]string) (string, error) {

	queryrow := `SELECT longurl from shorturls WHERE keyurl = $1`
	row := r.connPool.QueryRow(context.Background(), queryrow, key)
	var res string
	return res, row.Scan(&res)
}

// Get stored URL value by Key w/o username
func (r *repository) GetIDByLong(ctx context.Context, longURL, User string, UserList *[]string) (string, error) {
	queryrow := `SELECT keyurl from shorturls WHERE longurl = $1`
	row := r.connPool.QueryRow(context.Background(), queryrow, longURL)
	var res string
	return res, row.Scan(&res)
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
	return AllKeys, nil
}

func (r *repository) SetBatch(ctx context.Context, SetsForDB BatchSetsForDB) error {

	for _, Pair := range SetsForDB.Pairs {
		sqlInsert := "INSERT INTO shorturls (userid , keyurl , longurl) VALUES ($1, $2, $3)"
		_, err := r.connPool.Exec(ctx, sqlInsert, SetsForDB.UserID, Pair.ShortURL, Pair.LongURL)
		if err != nil {
			return errors.New("no way to implement this")
		}
	}
	return nil
}

func Contains(s *[]string, searchterm string) bool {
	i := sort.SearchStrings(*s, searchterm)
	return i < len(*s) && (*s)[i] == searchterm
}
