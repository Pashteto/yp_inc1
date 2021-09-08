package handlers

import (
	"context"
	"encoding/json"
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

	errReadDB := h.pingRedisDB(&h.Rdb)
	if errReadDB != nil {
		log.Println(errReadDB)
		http.Error(w, "DB not resonding", http.StatusInternalServerError)
		return
	}
	countID, _ := h.Rdb.Exists(ctx, id).Result()
	if countID == 0 {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, fmt.Sprintf("Wrong short URL id: %v", id), http.StatusBadRequest)
		return
	}
	longURL, _ := h.Rdb.Get(ctx, id).Result()
	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
	//	w.Write([]byte("Redirect"))
}

// Post puts the new url in the storage
func (h *HandlersWithDBStore) PostHandler(w http.ResponseWriter, r *http.Request) {
	errReadDB := h.pingRedisDB(&h.Rdb)
	if errReadDB != nil {
		log.Println(errReadDB)
		http.Error(w, "DB not resonding", http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "unable to parse body", http.StatusBadRequest)
		return
	}
	strbody := string(body)
	longURL, err := url.Parse(strbody)
	if err != nil {
		http.Error(w, "Unable to parse URL", http.StatusBadRequest)
		return
	}
	if len(strbody) == 0 {
		http.Error(w, "No URL recieved", http.StatusBadRequest)
		return
	}
	if !longURL.IsAbs() {
		longURL.Scheme = "http"
	}
	w.WriteHeader(http.StatusCreated)
	id := fmt.Sprint((rand.Intn(1000)))
	h.Rdb.Set(ctx, id, longURL.String(), 1000*time.Second)
	shorturl := config.String(h.Conf) + "/" + id
	w.Write([]byte(shorturl))
}

// Post puts the new url in the storage
func (h *HandlersWithDBStore) PostHandlerJSON(w http.ResponseWriter, r *http.Request) {
	errReadDB := h.pingRedisDB(&h.Rdb)
	if errReadDB != nil {
		log.Println(errReadDB)
		http.Error(w, "DB not resonding", http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "unable to parse body", http.StatusBadRequest)
		return
	}
	inputURL := fullURL{} //url.URL{} ///

	err = json.Unmarshal(body, &inputURL)
	if err != nil {
		log.Println(err)
		http.Error(w, "unable to unmarshal JSON", http.StatusBadRequest)
		return
	}
	if len(inputURL.CollectedString) == 0 {
		http.Error(w, "No URL recieved", http.StatusBadRequest)
		return
	}
	if !inputURL.collectedURL.IsAbs() {
		inputURL.collectedURL.Scheme = "http"
	}
	w.WriteHeader(http.StatusCreated)
	id := fmt.Sprint((rand.Intn(1000)))
	h.Rdb.Set(ctx, id, inputURL.collectedURL.String(), 1000*time.Second)
	outputURL := shortenResponse{}
	//outputURL.producedURL, err = url.Parse( config.String(h.Conf) + "/" + id)
	outputURL.ProducedString = config.String(h.Conf) + "/" + id

	output, err := json.Marshal(outputURL)
	if err != nil {
		http.Error(w, "unable to marshall short URL", http.StatusServiceUnavailable)
		return
	}

	w.Write(output)
}

func (h *HandlersWithDBStore) EmptyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func (h *HandlersWithDBStore) pingRedisDB(client *redis.Client) error {
	if client == nil {
		return errors.New("no redis db")
	}
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

type shortenResponse struct {
	ProducedString string `json:"result"`
	//	producedURL    *url.URL
}

type fullURL struct {
	CollectedString string `json:"url"`
	collectedURL    *url.URL
}

func (t *fullURL) UnmarshalJSON(data []byte) error {
	type fullURLAlias struct {
		CollectedString string `json:"url"`
	}
	aliasValue := fullURLAlias{}
	if err := json.Unmarshal(data, &aliasValue); err != nil {
		return err
	}
	var err error
	t.CollectedString = aliasValue.CollectedString
	t.collectedURL, err = url.Parse(aliasValue.CollectedString)
	if err != nil {
		return err
	}
	return nil
}
