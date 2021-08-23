package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
)

/*
Напишите сервис для сокращения длинных URL.
Требования:
Сервер должен быть доступен по адресу: http://localhost:8080.
Сервер должен предоставлять два эндпоинта: POST / и GET /{id}.
Эндпоинт POST / принимает в теле запроса строку URL для сокращения и возвращает в ответ правильный сокращённый URL.
Эндпоинт GET /{id} принимает в качестве URL параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400.
*/

func main() {

	// маршрутизация запросов обработчику
	http.HandleFunc("/", HelloServer)

	// конструируем свой сервер
	server := &http.Server{
		Addr: ":8080",
	}
	server.ListenAndServe()

	/*	// маршрутизация запросов обработчику
		http.HandleFunc("/", HelloServer)

		http.ListenAndServe(":8080", nil)
	*/

	// создаём канал для перехвата сигналов OS
	sigint := make(chan os.Signal, 1)
	// перенаправляем сигналы OS в этот канал
	signal.Notify(sigint, os.Interrupt)
	// ожидаем сигнала
	<-sigint
	// получаем сигнал OS и начинаем процедуру «мягкого останова»
	server.Shutdown(context.Background())
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!\n", r.URL.Path[1:])

	data := url.Values{}
	data.Set("id1", "google1.com")
	data.Set("id2", "google2.com")
	data.Set("id3", "google3.com")

	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("<h1>M Get</h1>"))
	case http.MethodPost:
		w.Write([]byte("<h1>M Post</h1>"))
		defer r.Body.Close()
		// читаем поток из тела ответа
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}
		// и печатаем его
		fmt.Println(string(body))
		fmt.Fprintf(w, "Hello, %v!\n", body)
		//prelim:=string(body)
		//prelim1:=strings.ReplaceAll(prelim,"url=","")
		//var id []string ={"id" "sd"}
		//numi:=string(32)
		//idnum := strings.Join(id,numi)

		//data.Set("id3", prelim1)
	default:
		http.Error(w, "Method not allowed, bad request", http.StatusBadRequest)
	}

}
