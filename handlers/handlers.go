package handlers

import (
	"context"
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

func (h *HandlersWithDBStore) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

// Handler for most of the bad requests
func (h *HandlersWithDBStore) GetHandler(w http.ResponseWriter, r *http.Request) {

	//id := mux.Vars(r)["key"]
	id := string(r.URL.Path[1:])
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

	var shorturl string
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	strbody := string(body)
	if err != nil {
		fmt.Println(err)
	} else {

		if !strings.Contains(strbody, "http:") && !strings.Contains(strbody, "https:") {
			strbody = "http://" + strbody
		}
		id := fmt.Sprint((rand.Intn(1000)))
		shorturl = "http://localhost/" + id
		h.Rdb.Set(ctx, id, strbody, 1000*time.Second)
	}
	fmt.Fprintf(w, "%v", shorturl)
}
