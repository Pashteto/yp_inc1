package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestHandlersWithDBStore_GetHandler(t *testing.T) {
	type fields struct {
		rdb            redis.Client
		code           int
		headerLocation string
		contentType    string
		id             string
		method         string
		///	response       string
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()

	rdb.Set(ctx, "this_id_is_a_correct_id", "http://google.com", 0)

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
			},
		},
		/*

			//      testing the DB: if it is initialised / connected?
			{
				name: "Test 3: Get Handler with no DB",
				fields: fields{
					code:           400,
					id:             "any_id",
					contentType:    "text/plain; charset=utf-8",
					headerLocation: "",
				},
			},*/

		//eo tests setup
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			hObj := &HandlersWithDBStore{
				Rdb: tt.fields.rdb,
			}
			request := httptest.NewRequest(tt.fields.method, "/"+tt.fields.id, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(hObj.GetHandler)
			h.ServeHTTP(w, request)
			res := w.Result()
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
		rdb          redis.Client
		code         int
		post_address string
		response     string
		method       string
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()

	tests := []struct {
		name   string
		fields fields
	}{
		// tests list
		{
			name: "Test 1: Post Handler fine response",
			fields: fields{
				rdb:  *rdb,
				code: 200,
				//				id:             "this_id_is_a_wrong_id",
				//contentType:    "text/plain; charset=utf-8",
				//headerLocation: "",
				post_address: "http://example.com",
				response:     "http://localhost/81",
				method:       "POST",
			},
		},
		{
			name: "Test 2: Post Handler empty body",
			fields: fields{
				rdb:  *rdb,
				code: 400,
				//				contentType:    "text/plain; charset=utf-8",
				//				headerLocation: "",
				post_address: "",
				response:     "No URL recieved\n",
				method:       "POST",
			},
		},
		/*
			{
			name: "Test 3: Post Handler DB broken",
			fields: fields{
				rdb:  *rdb,
				code: 400,
				//				contentType:    "text/plain; charset=utf-8",
				//				headerLocation: "",
				post_address: "",
				response:     "ERROR DB method\n",
				method: "",
			},
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hObj := &HandlersWithDBStore{
				Rdb: tt.fields.rdb,
			}

			endpoint := "http://localhost:8080/"
			data := tt.fields.post_address //url.Values{"url": {tt.fields.post_address}}

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
