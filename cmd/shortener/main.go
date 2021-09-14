package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Pashteto/yp_inc1/config"
	"github.com/Pashteto/yp_inc1/handlers"
	"github.com/Pashteto/yp_inc1/repos"

	"github.com/caarlos0/env/v6"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	//"github.com/joho/godotenv"
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
	/*
		fmt.Println(os.Getenv("USER"), "<<== that was username\n",
			os.ExpandEnv("home env is: ${HOME}\n"),
			os.ExpandEnv("port env is: ${PORT}\n"),
			os.ExpandEnv("files env is: ${FILES}\n"))*/

	var conf config.Config
	//config.ReadFile(&conf)

	//	os.Setenv("SERVER_HOST", "localhost")
	//	os.Setenv("PORT", "8080")
	//	os.Setenv("HOSTS", "localhost")

	err := env.Parse(&conf)
	conf.CheckEnv()
	if err != nil {
		log.Fatalf("Unable to read env:\t%v", err)
		//	t.Errorf("Unable to read config file conf.json:\t%v", err)
	}
	//	fmt.Printf("Current user is: %v\n", conf.User)
	fmt.Printf("CONF: %+v", conf)

	//godotenv.Load()
	/*	log.Println(os.Getenv("REDIS_HOST"),
			os.Getenv("APP_BASE_HOST"),
			os.Getenv("APP_PORT"),
			os.Getenv("APP_BASE_URL"))
		// conf.RecieveEnv(os.Getenv("APP_BASE_HOST"),
			// os.Getenv("APP_PORT"),
			// os.Getenv("APP_BASE_URL"))
	*/
	// initialising redis DB
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	repa := repos.NewRedisRepository(rdb)
	defer rdb.Close()

	// Passing the DB to the new obj with Handlers as methods
	sshand := handlers.HandlersWithDBStore{Rdb: &repa, Conf: &conf}

	r := mux.NewRouter()
	r.HandleFunc("/{key}", sshand.GetHandler).Methods("GET")             //routing get with the {key}
	r.HandleFunc("/api/shorten", sshand.PostHandlerJSON).Methods("POST") //routing post w JSON
	r.HandleFunc("/", sshand.PostHandler).Methods("POST")                //routing post

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
