package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/Pashteto/yp_inc1/handlers"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

var ctx = context.Background()

func main() {
	// initialising redis DB
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()
	// Passing the DB to the new obj with methods
	sshand := handlers.HandlersWithDBStore{Rdb: *rdb}

	r := mux.NewRouter()
	r.HandleFunc("/{key}", sshand.GetHandler).Methods("GET")
	r.HandleFunc("/", sshand.PostHandler).Methods("POST")
	r.HandleFunc("/", sshand.EmptyHandler)

	http.Handle("/", r)

	// конструируем свой сервер
	server := &http.Server{
		Addr: ":8080",
	}
	server.ListenAndServe()

	// создаём канал для перехвата сигналов OS bbb 	// перенаправляем сигналы OS в этот канал
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	// ожидаем сигнала
	<-sigint
	// получаем сигнал OS и начинаем процедуру «мягкого останова»
	server.Shutdown(ctx)
}
