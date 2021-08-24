package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

var ctx = context.Background()

//var dataB = url.Values{}

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func GetHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["key"]
	//	fmt.Println(maped)
	booid, _ := rdb.Exists(ctx, id).Result()
	// fmt.Println("key2", did)
	if booid > 0 {
		long_url, _ := rdb.Get(ctx, id).Result()
		http.Redirect(w, r, long_url, http.StatusTemporaryRedirect)
	} else {
		http.NotFound(w, r)
	}

}

func EmptyGetHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not allowed", http.StatusBadRequest)
}
func EmptyPostHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not allowed", http.StatusBadRequest)
}
func EmptyHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Add("google.com")
	http.Redirect(w, r, "https://google.com", http.StatusTemporaryRedirect)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {

	var shorturl string
	defer r.Body.Close()
	// читаем поток из тела ответа
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	} else {
		id := fmt.Sprint((rand.Intn(1000)))
		shorturl = "http://localhost/" + id
		rdb.Set(ctx, id, string(body), 10*time.Second)
	}
	fmt.Fprintf(w, "%v", shorturl)
}

func main() {

	r := mux.NewRouter()
	//	r.HandleFunc("/{key}/", GetHandler).Methods("GET")
	r.HandleFunc("/{key}", GetHandler).Methods("GET")
	r.HandleFunc("/", PostHandler).Methods("POST")
	r.HandleFunc("/", EmptyHandler)

	http.Handle("/", r)

	// конструируем свой сервер
	server := &http.Server{
		Addr: ":8080",
	}
	server.ListenAndServe()

	// создаём канал для перехвата сигналов OS
	sigint := make(chan os.Signal, 1)
	// перенаправляем сигналы OS в этот канал
	signal.Notify(sigint, os.Interrupt)
	// ожидаем сигнала
	<-sigint
	// получаем сигнал OS и начинаем процедуру «мягкого останова»
	server.Shutdown(context.Background())
}
