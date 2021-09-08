package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"

	"github.com/Pashteto/yp_inc1/config"
)

func TestHandlersWithDBStore_GetHandler(t *testing.T) {
	type fields struct {
		rdb            redis.Client
		code           int
		headerLocation string
		contentType    string
		id             string
		method         string
		conf           *config.Config
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()

	rdb.Set(ctx, "this_id_is_a_correct_id", "http://google.com", 0)

	var conf config.Config
	err := config.ReadFile(&conf)
	if err != nil {
		t.Errorf("Unable to read config file conf.json:\t%v", err)
		return
	}

	tests := []struct {
		name   string
		fields fields
	}{

		// tests list
		{
			name: "Test 1: Get Handler with wrong id",
			fields: fields{
				rdb:            *rdb,
				code:           400,
				id:             "this_id_is_a_wrong_id",
				contentType:    "text/plain; charset=utf-8",
				headerLocation: "",
				method:         "GET",
				conf:           &conf,
			},
		},
		{
			name: "Test 2: Get Handler with correct id",
			fields: fields{
				rdb:            *rdb,
				code:           307,
				contentType:    "text/html; charset=utf-8",
				id:             "this_id_is_a_correct_id",
				headerLocation: "http://google.com",
				method:         "GET",
				conf:           &conf,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			hObj := &HandlersWithDBStore{
				Rdb:  tt.fields.rdb,
				Conf: &conf,
			}
			request := httptest.NewRequest(tt.fields.method, "/"+tt.fields.id, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(hObj.GetHandler)
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			// test StatusCode
			if res.StatusCode != tt.fields.code {
				t.Errorf("Expected status code %d, got %d", tt.fields.code, w.Code)
			}
			// test Header Location
			if res.Header.Get("Location") != tt.fields.headerLocation {
				t.Errorf("Expected Location %s, got %s", tt.fields.headerLocation, res.Header.Get("Location"))
			}
			// test Header Content-Type
			assert.Equal(t, tt.fields.contentType, res.Header.Get("Content-Type"),
				"Expected Content-Type %s, got %s", tt.fields.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestHandlersWithDBStore_PostHandler(t *testing.T) {
	type fields struct {
		rdb         redis.Client
		code        int
		postAddress string
		response    string
		method      string
		conf        *config.Config
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()
	var conf config.Config
	err := config.ReadFile(&conf)
	if err != nil {
		t.Errorf("Unable to read config file conf.json:\t%v", err)
		return
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// tests list
		{
			name: "Test 1: Post Handler correct response",
			fields: fields{
				rdb:         *rdb,
				code:        201,
				postAddress: "http://example.com",
				response:    "http://localhost:8080/81",
				method:      "POST",
				conf:        &conf,
			},
		},
		{
			name: "Test 2: Post Handler empty body",
			fields: fields{
				rdb:         *rdb,
				code:        400,
				postAddress: "",
				response:    "No URL recieved\n",
				method:      "POST",
				conf:        &conf,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hObj := &HandlersWithDBStore{
				Rdb:  tt.fields.rdb,
				Conf: &conf,
			}

			endpoint := "http://localhost:8080/"
			data := tt.fields.postAddress

			request, err := http.NewRequest(tt.fields.method, endpoint, bytes.NewBufferString(data))
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()
			h := http.HandlerFunc(hObj.PostHandler)
			h.ServeHTTP(w, request)
			res := w.Result()

			// test StatusCode
			assert.Equal(t, tt.fields.code, res.StatusCode,
				"Expected status code %d, got %d", tt.fields.code, w.Code)

			// reading body
			defer res.Body.Close()
			resBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(resBody) != tt.fields.response {
				t.Errorf("Expected body %s, got %s", tt.fields.response, w.Body.String())
			}
		})
	}
}
