package middlewares

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

type Middleware func(http.Handler) http.Handler

func Сonveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

// middleware-функция принимает параметром Handler и возвращает тоже Handler.
func GzipMiddlewareRead(next http.Handler) http.Handler {
	// собираем Handler приведением типа
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()

		switch r.Header.Get("Content-Encoding") {
		case "gzip":
			//			var buf bytes.Buffer
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Printf("failed decompress data: %v", err)
			}
			defer reader.Close()
			//			io.Copy(&buf, reader) // print html to standard out
			//			r.Body = ioutil.NopCloser(reader)///&buf)
			r.Body = reader
			//			next.ServeHTTP(w, r)
		default:
		}
		next.ServeHTTP(w, r)
	})
}

func GzipMiddlewareWrite(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

/*
type gzipReader struct {
	http.Request
	Reader io.Reader
}

func (r *gzipReader) Read(b []byte) (int, error) {
	// Reader будет отвечать за gzip-чтение, поэтому пишем в него

	return r.Reader.Read(b)
}

func Decompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	defer r.Close()
	if err != nil {
		r.Close()
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}
	var b bytes.Buffer
	_, err = b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}
	return b.Bytes(), nil
}
*/
