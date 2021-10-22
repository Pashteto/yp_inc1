package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
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
	longURL, _ := h.Rdb.GetValueByKey(ctx, id, UserID.Value, &h.UsersInDB)
	if longURL == "" {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, fmt.Sprintf("Wrong short URL id: %v", id), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
}

// Get Handler provides with initial URLs stored by their ids
func (h *HandlersWithDBStore) GetPostgresPingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/sql")
	err := h.Rdb.Ping(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
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
		http.Error(w, "no URLs", http.StatusNoContent)
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

	var id string
	id, err = h.PostInDBReturnID(longURL, UserID.Value)

	if err != nil {
		errExists := errors.New(`longURL exists`)
		if err.Error() != errExists.Error() {
			http.Error(w, "No URL recieved", http.StatusBadRequest)
			return
		}
		id, _ = h.Rdb.GetIDByLong(ctx, longURL.String(), UserID.Value, &h.UsersInDB)
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
		filedb.WriteAll(h.Rdb, *h.Conf, &h.UsersInDB)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(h.Conf.BaseURL + "/" + id))

}

func (h *HandlersWithDBStore) PostInDBReturnID(longURL *url.URL, UserID string) (string, error) {
	if err := longURLcorrector(longURL); err != nil {
		return "", err
	}
	id := h.idSearch(UserID)
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
	w.Header().Set("Content-Type", "application/json")
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
	var id string
	id, err = h.PostInDBReturnID(inputURL.CollectedURL, UserID.Value)
	if err != nil {
		errExists := errors.New(`longURL exists`)
		if err.Error() != errExists.Error() {
			http.Error(w, "No URL recieved", http.StatusBadRequest)
			return
		}

		id, _ = h.Rdb.GetIDByLong(ctx, inputURL.CollectedURL.String(), UserID.Value, &h.UsersInDB)
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
		filedb.WriteAll(h.Rdb, *h.Conf, &h.UsersInDB)
	}
	outputURL := typeHandlingURL{}
	outputURL.CollectedURL, _ = url.Parse(h.Conf.BaseURL + "/" + id)
	output, err2 := json.Marshal(outputURL)
	if err2 != nil {
		log.Println(err2)
		http.Error(w, "unable to marshall short URL", http.StatusInternalServerError)
		return
	}
	w.Write(output)
}

// Post puts the new url in the storage with JSON input
func (h *HandlersWithDBStore) PostBatchHandler(w http.ResponseWriter, r *http.Request) {
	UserID, _ := r.Cookie(cookieName)
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read request", http.StatusBadRequest)
		return
	}
	inputURL, err := readUBatchURLs(body)
	if err != nil {
		http.Error(w, "unable to unmarshal JSON", http.StatusBadRequest)
		return
	}
	outputURL, setsForDB := h.convertBatchURLs(inputURL, UserID.Value)
	i := len(outputURL)
	j := len(inputURL)
	if i == 0 {
		http.Error(w, "Error reading URLs", http.StatusBadRequest)
		return
	}
	if i != j {
		log.Println("Summ of ", j-i, " URLs were dropped")
	}
	///
	err = h.Rdb.SetBatch(ctx, setsForDB)
	if err != nil {
		http.Error(w, "Error writing in DB", http.StatusBadRequest)
		return
	}
	output, err := json.MarshalIndent(outputURL, "", "  ")
	if err != nil {
		http.Error(w, "Error marshalling in responce", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(output)
	filedb.WriteAll(h.Rdb, *h.Conf, &h.UsersInDB)
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
type batchURLsID struct {
	ID      string `json:"correlation_id"`
	LongURL string `json:"original_url"`
}
type batchShortURLsID struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

func readUBatchURLs(body []byte) ([]batchURLsID, error) {
	URLsID := []batchURLsID{}
	buf := bytes.NewBuffer(body)
	encoder := json.NewDecoder(buf)
	//encoder.SetEscapeHTML(false)
	err := encoder.Decode(&URLsID)
	if err != nil {
		return nil, err
	}
	return URLsID, nil
}

func (h *HandlersWithDBStore) convertBatchURLs(batchlongURLs []batchURLsID, UserID string) ([]batchShortURLsID, repos.BatchSetsForDB) {
	var shortURLsIDs []batchShortURLsID
	setsDB := repos.BatchSetsForDB{UserID: UserID, Pairs: []repos.IDShortURL{}}
	for _, batchlongURL := range batchlongURLs {
		longURL, err := url.Parse(batchlongURL.LongURL)
		if err != nil {
			continue
		}
		if err = longURLcorrector(longURL); err != nil {
			continue
		}
		id := h.idSearch(UserID)
		ShortURL := id //h.Conf.BaseURL + "/" + id
		shortURLsID := batchShortURLsID{ID: batchlongURL.ID, ShortURL: h.Conf.BaseURL + "/" + id}
		shortURLsIDs = append(shortURLsIDs, shortURLsID)

		addInDB := repos.IDShortURL{ShortURL: ShortURL, LongURL: longURL.String()}
		setsDB.Pairs = append(setsDB.Pairs, addInDB)
	}
	return shortURLsIDs, setsDB
}

func longURLcorrector(longURL *url.URL) error {
	if longURL.Host == "" && longURL.Path == "" {
		return errors.New("no URL recieved")
	}
	if !longURL.IsAbs() {
		longURL.Scheme = "http"
	}
	return nil
}

func (h *HandlersWithDBStore) idSearch(UserID string) string {
	var id string
	for {
		id = fmt.Sprint((rand.Intn(1000)))
		voidURL2, _ := h.Rdb.GetValueByKey(ctx, id, UserID, &h.UsersInDB)
		if voidURL2 == "" {
			break
		}
	}
	return id
}
