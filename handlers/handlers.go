package handlers

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx, cancel = context.WithCancel(context.Background())

// Storing data in this structure to get rid of global var DB
// data is stored using Redis DB
type HandlersWithDBStore struct {
	Rdb redis.Client
}

//func (h *HandlersWithDBStore) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

// Handler for most of the bad requests
func (h *HandlersWithDBStore) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Get handler recieved wrong method", http.StatusBadRequest)
		return
	}

	id := string(r.URL.Path[1:])

	errRedDb := pingRedisDb(&h.Rdb)
	if errRedDb != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "DB not resonding", http.StatusBadRequest)
		return
	}
	booid, _ := h.Rdb.Exists(ctx, id).Result()
	if booid > 0 {
		long_url, _ := h.Rdb.Get(ctx, id).Result()
		w.Header().Set("Content-Type", "text/html")
		http.Redirect(w, r, long_url, http.StatusTemporaryRedirect)
		w.Write([]byte("Redirect"))

	} else {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "Wrong id", http.StatusBadRequest)
	}

}

// Handler for most of the bad requests
func (h *HandlersWithDBStore) EmptyHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not allowed", http.StatusBadRequest)
}

// Puts the new url in the storage
func (h *HandlersWithDBStore) PostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Post handler recieved wrong method", http.StatusBadRequest)
		return
	}
	var shorturl string
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	strbody := string(body)
	if err != nil {
		fmt.Println(err)
	} else {
		if len(strbody) > 0 {
			if !strings.Contains(strbody, "http:") && !strings.Contains(strbody, "https:") {
				strbody = "http://" + strbody
			}
			id := fmt.Sprint((rand.Intn(1000)))
			shorturl = "http://localhost/" + id
			h.Rdb.Set(ctx, id, strbody, 1000*time.Second)
		} else {
			http.Error(w, "No URL recieved", http.StatusBadRequest)
		}
	}
	//	fmt.Fprintf(w, "%v", shorturl)
	w.Write([]byte(shorturl))

}

func pingRedisDb(client *redis.Client) error {
	//	if client {

	//	}
	if client == nil {
		err := errors.New("No Redis DB")
		return err
	}
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	//fmt.Println(pong, err)
	return nil
}
