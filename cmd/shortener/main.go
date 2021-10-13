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
<<<<<<< HEAD
	middlewares "github.com/Pashteto/yp_inc1/mddlwrs"
	"github.com/Pashteto/yp_inc1/repos"
	"github.com/caarlos0/env/v6"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
=======
	"github.com/Pashteto/yp_inc1/repos"
	"github.com/caarlos0/env/v6"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
>>>>>>> main
)

var ctx = context.Background()

func main() {
	var conf config.Config
<<<<<<< HEAD
	defaultPsqlConn := "host=localhost port=5432 user=postgres password=kornkorn dbname=mydb sslmode=disable"
=======
>>>>>>> main

	ServAddrPtr := flag.String("a", ":8080", "SERVER_ADDRESS")
	BaseURLPtr := flag.String("b", "http://localhost:8080", "BASE_URL")
	FStorPathPtr := flag.String("f", "../URLs", "FILE_STORAGE_PATH")
<<<<<<< HEAD
	PostgresURL := flag.String("d", defaultPsqlConn, "DATABASE_URL")

	flag.Parse()

	log.Println("Flags input:\nSERVER_ADDRESS,\tBASE_URL,\tFILE_STORAGE_PATH:\t",
		*ServAddrPtr, ",", *BaseURLPtr, ",", *FStorPathPtr)
=======
	flag.Parse()

	log.Println("Flags input:\nSERVER_ADDRESS,\tBASE_URL,\tFILE_STORAGE_PATH:\t",
		*ServAddrPtr, ",",
		*BaseURLPtr, ",", *FStorPathPtr)
>>>>>>> main
	err := env.Parse(&conf)
	if err != nil {
		log.Fatalf("Unable to Parse env:\t%v", err)
	}
<<<<<<< HEAD
	changed, err := conf.UpdateByFlags(ServAddrPtr, BaseURLPtr, FStorPathPtr, PostgresURL)
	if changed {
		log.Printf("Config updated:SERVER_ADDRESS:\t%v,BASE_URL:\t%v,FILE_STORAGE_PATH:\t%v,\n",
			conf.ServAddr, conf.BaseURL, conf.FStorPath)
=======
	log.Printf("Config:\t%+v", conf)

	changed, err := conf.UpdateByFlags(ServAddrPtr, BaseURLPtr, FStorPathPtr)
	if changed {
		log.Printf("Config updated:\t%+v\n", conf)
>>>>>>> main
	}
	if err != nil {
		log.Printf("Flags input error:\t%v\n", err)
	}

<<<<<<< HEAD
	log.Println("USER:\t", os.Getenv("USER"))
	// REDIS
	/*
		// initialising redis DB
		rdb := redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_HOST") + ":6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		defer rdb.Close()
	*/

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
=======
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
		log.Fatalf("file exited;\nerr:\t%v; repository:\t%v", err, repa)
	}
	err = filedb.UpdateDBSlice(repa, conf)
	if err != nil {
		log.Fatal(err)
	}
	defer rdb.Close()
	sshand := handlers.HandlersWithDBStore{Rdb: repa, Conf: &conf}

	r := mux.NewRouter()
	r.HandleFunc("/{key}", sshand.GetHandler).Methods("GET")             //routing get with the {key}
	r.HandleFunc("/api/shorten", sshand.PostHandlerJSON).Methods("POST") //routing post w JSON
	r.HandleFunc("/", sshand.PostHandler).Methods("POST")                //routing post
>>>>>>> main

	http.Handle("/", r)

	// конструируем свой сервер
	server := &http.Server{
<<<<<<< HEAD
		Addr: conf.ServAddr,
	}
=======

		Addr: conf.ServAddr,
	}

>>>>>>> main
	sigint := make(chan os.Signal)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	go func() {
		sig := <-sigint // Blocks here until interrupted
		log.Println(sig, "\t<<<===\t signal received. Shutdown process initiated.")
		server.Shutdown(ctx)
<<<<<<< HEAD
		filedb.WriteAll(sshand.Rdb, *sshand.Conf, &sshand.UsersInDB)
	}()
=======
	}()

>>>>>>> main
	server.ListenAndServe()
}
