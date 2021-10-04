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
	"sort"
	"time"

	"github.com/Pashteto/yp_inc1/config"
	filedb "github.com/Pashteto/yp_inc1/filed_history"
	"github.com/Pashteto/yp_inc1/repos"
)

var ctx, _ = context.WithCancel(context.Background())

const urlTTL = time.Second * 1000
const cookieName = "UserID"

// Storing data in this structure to get rid of global var DB
// data is stored using Redis DB
type HandlersWithDBStore struct {
	Rdb       repos.SetterGetter // redis.Client
	Conf      *config.Config
	UsersInDB []string
}

// Get Handler provides with initial URLs stored by their ids
func (h *HandlersWithDBStore) GetHandler(w http.ResponseWriter, r *http.Request) {
	id := string(r.URL.Path[1:])
	//	Checked r.cookie in middleware. It should be there
	UserID, _ := r.Cookie(cookieName)
	longURL1, _ := h.Rdb.GetValueByKey(ctx, id, UserID.Value, &h.UsersInDB)
	if longURL1 == "" {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, fmt.Sprintf("Wrong short URL id: %v", id), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, longURL1, http.StatusTemporaryRedirect)
}

// Get All Urls Handler provides with all URLs stored by single user
func (h *HandlersWithDBStore) GetAllUrlsHandler(w http.ResponseWriter, r *http.Request) {
	UserID, _ := r.Cookie(cookieName)
	//fmt.Printf("fi cookies are:\t%v", UserID)
	keys, err := h.Rdb.ListAllKeysByUser(ctx, UserID.Value)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, fmt.Sprintf("error receiving all key-value pairs from redis db: %v", err), http.StatusBadRequest)
		return
	}
	rf := len(keys)
	if rf == 0 {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, fmt.Sprintf("no URLs"), http.StatusNoContent)
		return
	}
	sliceIDURL := make([]iDShortURLReflex, rf)
	i := 0
	for key, val := range keys {
		sliceIDURL[i] = iDShortURLReflex{ID: h.Conf.BaseURL + "/" + key, LongURL: val}
		i++
	}
	output, err2 := json.MarshalIndent(sliceIDURL, "", "    ")
	if err2 != nil {
		http.Error(w, "unable to marshall all the short urls list", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

// Post puts the new url in the storage
func (h *HandlersWithDBStore) PostHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "unable to parse body", http.StatusBadRequest)
		return
	}
	longURL, err := url.Parse(string(body))
	if err != nil {
		http.Error(w, "Unable to parse URL", http.StatusBadRequest)
		return
	}
	UserID, _ := r.Cookie(cookieName)
	id, err := h.PostInDBReturnID(longURL, UserID.Value)

	if err != nil {
		http.Error(w, "No URL recieved", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(h.Conf.BaseURL + "/" + id))
	filedb.WriteAll(h.Rdb, *h.Conf, &h.UsersInDB)
}

func (h *HandlersWithDBStore) PostInDBReturnID(longURL *url.URL, UserID string) (string, error) {
	if longURL.Host == "" && longURL.Path == "" {
		return "", errors.New("no URL recieved")
	}
	if !longURL.IsAbs() {
		longURL.Scheme = "http"
	}
	var id string
	for {
		id = fmt.Sprint((rand.Intn(1000)))
		voidURL2, _ := h.Rdb.GetValueByKey(ctx, id, UserID, &h.UsersInDB)
		if voidURL2 == "" {
			break
		}
	}
	err2 := h.Rdb.SetValueByKeyAndUser(ctx, id, UserID, longURL.String(), urlTTL)
	if err2 != nil {
		return id, err2
	}
	if !repos.Contains(&h.UsersInDB, UserID) {
		h.UsersInDB = append(h.UsersInDB, UserID)
		sort.Strings(h.UsersInDB)
	}
	return id, err2
}

// Post puts the new url in the storage with JSON input
func (h *HandlersWithDBStore) PostHandlerJSON(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "unable to parse body", http.StatusBadRequest)
		return
	}
	inputURL := typeHandlingURL{}
	err = json.Unmarshal(body, &inputURL)
	if err != nil {
		log.Println(err)
		http.Error(w, "unable to unmarshal JSON", http.StatusBadRequest)
		return
	}
	UserID, _ := r.Cookie(cookieName)
	id, err := h.PostInDBReturnID(inputURL.CollectedURL, UserID.Value)
	if err != nil {
		http.Error(w, "No URL recieved", http.StatusBadRequest)
		return
	}
	outputURL := typeHandlingURL{}
	outputURL.CollectedURL, _ = url.Parse(h.Conf.BaseURL + "/" + id)
	output, err2 := json.Marshal(outputURL)
	if err2 != nil {
		http.Error(w, "unable to marshall short URL", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(output)
	filedb.WriteAll(h.Rdb, *h.Conf, &h.UsersInDB)
}

func (h *HandlersWithDBStore) EmptyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

type typeHandlingURL struct {
	CollectedURL *url.URL
}

func (t typeHandlingURL) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ResultString string `json:"result"`
	}{
		ResultString: t.CollectedURL.String(),
	})
}

func (t *typeHandlingURL) UnmarshalJSON(data []byte) error {
	type typeHandlingURLAlias struct {
		CollectedString string `json:"url"`
	}
	aliasValue := typeHandlingURLAlias{}
	if err := json.Unmarshal(data, &aliasValue); err != nil {
		return err
	}
	var err error
	t.CollectedURL, err = url.Parse(aliasValue.CollectedString)
	if err != nil {
		return err
	}
	return nil
}

type iDShortURLReflex struct {
	ID      string `json:"short_url"`
	LongURL string `json:"original_url"`
}
