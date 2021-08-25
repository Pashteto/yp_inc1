package handlers

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

var ctx, cancel = context.WithCancel(context.Background())

type SpecificHandler struct {
	Rdb redis.Client
}

func (h *SpecificHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func (h *SpecificHandler) GetHandler(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["key"]
	//	fmt.Println(maped)
	booid, _ := h.Rdb.Exists(ctx, id).Result()
	// fmt.Println("key2", did)
	if booid > 0 {
		long_url, _ := h.Rdb.Get(ctx, id).Result()
		http.Redirect(w, r, long_url, http.StatusTemporaryRedirect)
	} else {
		http.NotFound(w, r)
	}

}

/*func (h *SpecificHandler) EmptyGetHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not allowed", http.StatusBadRequest)
}
func (h *SpecificHandler) EmptyPostHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not allowed", http.StatusBadRequest)
}*/
func (h *SpecificHandler) EmptyHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Add("google.com")
	//	http.Redirect(w, r, "https://google.com", http.StatusTemporaryRedirect)
	http.Error(w, "Method not allowed", http.StatusBadRequest)
}

func (h *SpecificHandler) PostHandler(w http.ResponseWriter, r *http.Request) {

	var shorturl string
	defer r.Body.Close()
	// читаем поток из тела ответа
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	} else {
		id := fmt.Sprint((rand.Intn(1000)))
		shorturl = "http://localhost/" + id
		h.Rdb.Set(ctx, id, string(body), 1000*time.Second)
	}
	fmt.Fprintf(w, "%v", shorturl)
}
