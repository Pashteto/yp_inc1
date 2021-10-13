package middlewares

import (
	"compress/gzip"
	"crypto/rand"
	"io"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"
)

const urlTTL = time.Second * 1000
const sizeOfKey = 15
const cookieName = "UserID"

type Middleware func(http.Handler) http.Handler

func Сonveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

// middleware-функция чтения архивированного тела запроса
func GzipMiddlewareRead(next http.Handler) http.Handler {
	// возвращаем Handler приведением типа
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()
		log.Println("GzipMiddlewareRead usage")

		switch r.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Printf("failed decompress data: %v", err)
			}
			defer reader.Close()
			r.Body = reader
		default:
		}
		next.ServeHTTP(w, r)
	})
}

// middleware-функция записи архивированного ответа
func GzipMiddlewareWrite(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			log.Println("GzipMiddlewareWrite midlware easy")

			return
		}
		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			log.Printf("failed to compress data: %v", err)
			next.ServeHTTP(w, r)
			log.Println("GzipMiddlewareWrite err")

			return
		}
		defer gz.Close()
		log.Println("GzipMiddlewareWrite next")

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Writer будет отвечать за gzip-сжатие, поэтому пишем в него
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// middleware-функция записи архивированного ответа
func UserCookieCheckGen(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, есть ли у запроса куки
		fi, err := r.Cookie(cookieName)
		if fi != nil {
			next.ServeHTTP(w, r)
			return
		}
		if err != nil && err.Error() != "http: named cookie not present" {
			log.Println("Request's cookie parsing error ocurred.")
			next.ServeHTTP(w, r)
			return
		}
		expiration := time.Now().Add(urlTTL)
		UserID, err := generateRandomString(sizeOfKey) // ключ шифрования
		if err != nil {
			log.Printf("error generating user ID: %v\n", err)
			UserID = "abcd"
		}
		cookie := http.Cookie{Name: cookieName, Value: UserID, Expires: expiration}
		http.SetCookie(w, &cookie)
		r.AddCookie(&cookie)
		next.ServeHTTP(w, r)
	})
}

func generateRandomString(size int) (string, error) {
	letters := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-+!@#%^*"
	lettersNum := len(letters)
	ret := make([]byte, size)
	for i := 0; i < size; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(lettersNum)))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}
	return string(ret), nil
}
