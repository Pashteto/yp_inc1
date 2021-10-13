package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Pashteto/yp_inc1/config"
	filedb "github.com/Pashteto/yp_inc1/filed_history"
	"github.com/Pashteto/yp_inc1/handlers"
	middlewares "github.com/Pashteto/yp_inc1/mddlwrs"
	"github.com/Pashteto/yp_inc1/repos"
	"github.com/caarlos0/env/v6"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ctx = context.Background()

func main() {
	var conf config.Config
	defaultPsqlConn := "host=localhost port=5432 user=postgres password=kornkorn dbname=mydb sslmode=disable"

	ServAddrPtr := flag.String("a", ":8080", "SERVER_ADDRESS")
	BaseURLPtr := flag.String("b", "http://localhost:8080", "BASE_URL")
	FStorPathPtr := flag.String("f", "../URLs", "FILE_STORAGE_PATH")
	PostgresURL := flag.String("d", defaultPsqlConn, "DATABASE_URL")

	flag.Parse()

	log.Println("Flags input:\nSERVER_ADDRESS,\tBASE_URL,\tFILE_STORAGE_PATH:\t",
		*ServAddrPtr, ",", *BaseURLPtr, ",", *FStorPathPtr)
	err := env.Parse(&conf)
	if err != nil {
		log.Fatalf("Unable to Parse env:\t%v", err)
	}
	changed, err := conf.UpdateByFlags(ServAddrPtr, BaseURLPtr, FStorPathPtr, PostgresURL)
	if changed {
		log.Printf("Config updated:SERVER_ADDRESS:\t%v,BASE_URL:\t%v,FILE_STORAGE_PATH:\t%v,\n",
			conf.ServAddr, conf.BaseURL, conf.FStorPath)
	}
	if err != nil {
		log.Printf("Flags input error:\t%v\n", err)
	}

	pool, err := pgxpool.Connect(context.Background(), conf.PostgresURL)
	if err != nil {
		log.Fatalf("Postgres connect error:\t%v", err)
	}
	defer pool.Close()
	repa, err := repos.NewRepoWitnTable(ctx, pool)
	if err != nil {
		log.Fatalf("Error creating repository:\terr:\t%v; repository:\t%v", err, repa)
	}
	err = filedb.CreateDirFileDBExists(conf)
	if err != nil {
		log.Fatalf("file exited;\nerr:\t%v", err)
	}

	UserList, err := filedb.UpdateDBSlice(repa, conf)
	if err != nil {
		log.Fatal(err)
	}
	sshand := handlers.HandlersWithDBStore{Rdb: repa, Conf: &conf, UsersInDB: UserList}
	r := mux.NewRouter()
	r.HandleFunc("/user/urls", sshand.GetAllUrlsHandler).Methods("GET") //routing get for all the keys of this user
	r.HandleFunc("/ping", sshand.GetPostgresPingHandler).Methods("GET") //routing ping of the postgres db
	r.HandleFunc("/{key}", sshand.GetHandler).Methods("GET")            //routing get with the {key}

	r.HandleFunc("/api/shorten", sshand.PostHandlerJSON).Methods("POST") //routing post w JSON
	r.HandleFunc("/", sshand.PostHandler).Methods("POST")                //routing post
	r.Use(middlewares.UserCookieCheckGen)
	r.Use(middlewares.GzipMiddlewareRead)
	r.Use(middlewares.GzipMiddlewareWrite)

	http.Handle("/", r)

	// конструируем свой сервер
	server := &http.Server{
		Addr: conf.ServAddr,
	}
	sigint := make(chan os.Signal)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	go func() {
		sig := <-sigint // Blocks here until interrupted
		log.Println(sig, "\t<<<===\t signal received. Shutdown process initiated.")
		server.Shutdown(ctx)
		filedb.WriteAll(sshand.Rdb, *sshand.Conf, &sshand.UsersInDB)
	}()
	server.ListenAndServe()
}
