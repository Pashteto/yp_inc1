package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Pashteto/yp_inc1/config"
	"github.com/Pashteto/yp_inc1/handlers"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

/*
Задание для трека «Go в веб-разработке». Напишите сервис для сокращения длинных URL. Требования:
- Сервер должен быть доступен по адресу: http://localhost:8080.
- Сервер должен предоставлять два эндпоинта: POST / и GET /{id}.
- Эндпоинт POST / принимает в теле запроса строку URL для сокращения и возвращает в ответ правильный сокращённый URL.
- Эндпоинт GET /{id} принимает в качестве URL параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
- Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400.
*/

var ctx = context.Background()

func main() {
	var conf config.Config
	err := config.ReadFile(&conf)
	if err != nil {
		log.Println("Unable to read config file conf.json:\t", err)
	}

	err2 := godotenv.Load()
	if err2 != nil {
		log.Fatal("Error loading .env file")
	}

	// initialising redis DB
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()

	// Passing the DB to the new obj with Handlers as methods
	sshand := handlers.HandlersWithDBStore{Rdb: *rdb, Conf: &conf}

	r := mux.NewRouter()
	r.HandleFunc("/{key}", sshand.GetHandler).Methods("GET") //routing get with the {key}
	r.HandleFunc("/", sshand.PostHandler).Methods("POST")    //routing post
	r.HandleFunc("/", sshand.EmptyHandler)                   //routing post

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
	// получаем сигнал OS и начинаем процедуру «мягкой остановки»
	server.Shutdown(ctx)
}
