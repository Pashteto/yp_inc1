package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

var ctx = context.Background()
var dataB = url.Values{}

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func GetHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "GetHandler, %v!\n", mux.Vars(r))
	//fmt.Fprintf(w, "url data, %v!\n", r.URL.Path[1:])

	//fmt.Println(mux.Vars(r))
	//
	//	id := r.URL.Path[1:]
	//	var long_url string
	defer r.Body.Close()
	// читаем поток из тела ответа
	/* body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	} else {
	*/
	/*

	   	did,_:=rdb.Exists(ctx, "key2").Result()
	   	fmt.Println("key2", did)
	   	if did > 0  {
	   		val2, _ := rdb.Get(ctx, "key2").Result()
	   	} else

	       val2, _ := rdb.Get(ctx, "key2").Result()
	       if err == redis.Nil {
	           fmt.Println("key2 does not exist")
	       } else if err != nil {
	           panic(err)
	       } else {
	           fmt.Println("key2", val2)
	       }

	   	if found_shorturl {

	   		http.Redirect(w, r, long_url, http.StatusTemporaryRedirect)
	   	} else {
	   		http.NotFound(w, r)
	   	}*/
	//fmt.Fprintf(w, "PostHandler, %v!\n", mux.Vars(r))
	//fmt.Fprintf(w, "url data, %v!\n", r.URL.Path[1:])
	//}
	//fmt.Println(string(body))
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
		addthis := true
		for it, it2 := range dataB {
			if it == string(body) {
				addthis = false
				shorturl = strings.Join(it2, "")
				break
			}
		}
		if addthis {
			shorturl = "http://localhost/" + fmt.Sprint((rand.Intn(1000)))

			dataB.Set(string(body), shorturl)

			fmt.Println(rdb)

			//	fmt.Println(shorturl)
		}
		//fmt.Fprintf(w, "PostHandler, %v!\n", mux.Vars(r))
		//fmt.Fprintf(w, "url data, %v!\n", r.URL.Path[1:])
	}
	//fmt.Println(string(body))

	fmt.Fprintf(w, "%v", shorturl)
}

func main() {

	r := mux.NewRouter()
	//	r.HandleFunc("/{key}", GetHandler).Methods("GET")
	//	r.HandleFunc("/", PostHandler).Methods("POST")
	r.HandleFunc("/", EmptyHandler)

	http.Handle("/", r)
	rdb.Set(ctx, "ASS", "Ssdkhb", 0)
	val2, _ := rdb.Get(ctx, "ASS").Result()
	fmt.Println(val2)
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
