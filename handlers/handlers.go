package handlers

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
)

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})
var ctxH, cancel = context.WithCancel(context.Background())

func GetHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	} else {
		did, _ := rdb.Exists(ctxH, string(body)).Result()
		fmt.Println("key2", did)
		if did > 0 {
			val2, _ := rdb.Get(ctxH, string(body)).Result()
			http.Redirect(w, r, val2, http.StatusTemporaryRedirect)
		} else {
			http.NotFound(w, r)
		}
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
