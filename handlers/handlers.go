package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/Pashteto/yp_inc1/config"
	"github.com/go-redis/redis/v8"
)

var ctx, _ = context.WithCancel(context.Background())

// Storing data in this structure to get rid of global var DB
// data is stored using Redis DB
type HandlersWithDBStore struct {
	Rdb  redis.Client
	Conf *config.Config
}

// Get Handler provides with initial URLs stored by their ids
func (h *HandlersWithDBStore) GetHandler(w http.ResponseWriter, r *http.Request) {
	id := string(r.URL.Path[1:])

	errReadDb := h.pingRedisDb(&h.Rdb)
	if errReadDb != nil {
		log.Println(errReadDb)
		http.Error(w, "DB not resonding", http.StatusInternalServerError)
		return
	}
	count_id, _ := h.Rdb.Exists(ctx, id).Result()
	if count_id == 0 {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, fmt.Sprintf("Wrong short URL id: %v", id), http.StatusBadRequest)
		return
	}
	long_url, _ := h.Rdb.Get(ctx, id).Result()
	http.Redirect(w, r, long_url, http.StatusTemporaryRedirect)
	w.Write([]byte("Redirect"))
}

// Post puts the new url in the storage
func (h *HandlersWithDBStore) PostHandler(w http.ResponseWriter, r *http.Request) {
	errReadDb := h.pingRedisDb(&h.Rdb)
	if errReadDb != nil {
		log.Println(errReadDb)
		http.Error(w, "DB not resonding", http.StatusInternalServerError)
		return
	}
	var shorturl string

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "unable to parse body", http.StatusBadRequest)
		return
	}
	strbody := string(body)
	longUrl, err := url.Parse(strbody)
	if err != nil {
		http.Error(w, "Unable to parse URL", http.StatusBadRequest)
		return
	}
	if len(strbody) == 0 {
		http.Error(w, "No URL recieved", http.StatusBadRequest)
		return
	}
	if !longUrl.IsAbs() {
		longUrl.Scheme = "http"
	}
	id := fmt.Sprint((rand.Intn(1000)))
	shorturl = config.String(h.Conf) + "/" + id
	h.Rdb.Set(ctx, id, longUrl.String(), 1000*time.Second)
	w.Write([]byte(shorturl))
}

func (h *HandlersWithDBStore) pingRedisDb(client *redis.Client) error {
	if client == nil {
		return errors.New("no redis db")
	}
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}
