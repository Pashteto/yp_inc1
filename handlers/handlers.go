package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx, _ = context.WithCancel(context.Background())

// Storing data in this structure to get rid of global var DB
// data is stored using Redis DB
type HandlersWithDBStore struct {
	Rdb redis.Client
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
	booid, _ := h.Rdb.Exists(ctx, id).Result()
	if booid > 0 {
		long_url, _ := h.Rdb.Get(ctx, id).Result()
		http.Redirect(w, r, long_url, http.StatusTemporaryRedirect)
		w.Write([]byte("Redirect"))

	} else {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, fmt.Sprintf("Wrong short URL id: %v", id), http.StatusBadRequest)
	}

}

/* Commented due to being unnesessary - the router does this automatically
// Handler for Bad requests
func (h *HandlersWithDBStore) HandlerBadRequest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Bad request", http.StatusBadRequest)
}*/

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
	if len(strbody) > 0 {
		if !strings.Contains(strbody, "http:") && !strings.Contains(strbody, "https:") {
			strbody = "http://" + strbody
		}
		id := fmt.Sprint((rand.Intn(1000)))
		hostProto := "http://"
		if r.TLS != nil {
			hostProto = "https://"
		}
		hostName := hostProto + strings.Split(r.Host, ":")[0]
		shorturl = hostName + "/" + id
		h.Rdb.Set(ctx, id, strbody, 1000*time.Second)
	} else {
		http.Error(w, "No URL recieved", http.StatusBadRequest)
	}

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

//func (h *HandlersWithDBStore) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
