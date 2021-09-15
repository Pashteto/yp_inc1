package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Pashteto/yp_inc1/config"
	filedb "github.com/Pashteto/yp_inc1/filed_history"
	"github.com/Pashteto/yp_inc1/handlers"
	"github.com/Pashteto/yp_inc1/repos"
	"github.com/caarlos0/env/v6"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

var ctx = context.Background()

func main() {
	var conf config.Config

	ServAddrPtr := flag.String("a", ":8080", "SERVER_ADDRESS")
	BaseURLPtr := flag.String("b", "http://localhost:8080", "BASE_URL")
	FStorPathPtr := flag.String("f", os.Getenv("HOME"), "FILE_STORAGE_PATH")
	RedisPtr := flag.String("r", os.Getenv("REDIS_HOST"), "REDIS_HOST")

	flag.Parse()
	log.Println("\nFlags input:\nSERVER_ADDRESS,\tBASE_URL,\tFILE_STORAGE_PATH:\n", *ServAddrPtr, ",", *BaseURLPtr, ",", *FStorPathPtr)
	err := env.Parse(&conf)
	if err != nil {
		log.Fatalf("Unable to Parse env:\t%v", err)
	}
	log.Printf("Config:\t%+v", conf)
	if conf.UpdateByFlags(ServAddrPtr, BaseURLPtr, FStorPathPtr, RedisPtr) {
		log.Printf("Config updated:\t%+v", conf)
	}

	log.Println("REDIS_HOST:\t", os.Getenv("REDIS_HOST"))
	log.Println("USER:\t", os.Getenv("USER"))

	// initialising redis DB
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	repa := repos.NewRedisRepository(rdb)
	err = filedb.CreateDirFileDBExists(conf)
	if err != nil {
		log.Printf("file exited")
		log.Fatal(err, repa)
	}
	//err = filedb.UpdateDB(&repa, conf)
	err = filedb.UpdateDBSlice(&repa, conf)

	if err != nil {
		log.Fatal(err)
	}
	defer rdb.Close()
	sshand := handlers.HandlersWithDBStore{Rdb: &repa, Conf: &conf}

	r := mux.NewRouter()
	r.HandleFunc("/{key}", sshand.GetHandler).Methods("GET")             //routing get with the {key}
	r.HandleFunc("/api/shorten", sshand.PostHandlerJSON).Methods("POST") //routing post w JSON
	r.HandleFunc("/", sshand.PostHandler).Methods("POST")                //routing post

	http.Handle("/", r)

	// конструируем свой сервер
	server := &http.Server{
		Addr: conf.ServAddr,
	}
	server.ListenAndServe()

	err = filedb.WriteAll(&repa, conf)
	// создаём канал для перехвата сигналов OS bbb 	// перенаправляем сигналы OS в этот канал
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	// ожидаем сигнала
	<-sigint
	// получаем сигнал OS и начинаем процедуру «мягкой остановки»
	server.Shutdown(ctx)

}
