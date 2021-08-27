package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-redis/redis/v8"
)

func TestHandlersWithDBStore_GetHandler(t *testing.T) {
	type fields struct {
		rdb            redis.Client
		code           int
		headerLocation string
		contentType    string
		id             string
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	rdb.Set(ctx, "this_id_is_a_correct_id", "http://google.com", 0)

	tests := []struct {
		name   string
		fields fields
		//args   args
	}{

		//tests setup
		{
			name: "Test1: Get Handler with wrong id",
			fields: fields{
				rdb:            *rdb,
				code:           400,
				id:             "this_id_is_a_wrong_id",
				contentType:    "text/plain; charset=utf-8",
				headerLocation: "",
			},
		},
		{
			name: "Test2: Get Handler with correct id",
			fields: fields{
				rdb:            *rdb,
				code:           307,
				contentType:    "text/html",
				id:             "this_id_is_a_correct_id",
				headerLocation: "http://google.com",
			},
		},

		//eo tests setup
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			hObj := &HandlersWithDBStore{
				Rdb: tt.fields.rdb,
			}

			request := httptest.NewRequest("GET", "/"+tt.fields.id, nil)
			//request = tt.args.r
			w := httptest.NewRecorder()
			h := http.HandlerFunc(hObj.GetHandler)
			h.ServeHTTP(w, request)
			//			hObj.GetHandler(tt.args.w, tt.args.r)
			res := w.Result()
			// проверяем код ответа
			if res.StatusCode != tt.fields.code {
				t.Errorf("Expected status code %d, got %d", tt.fields.code, w.Code)
			}
			if res.Header.Get("Location") != tt.fields.headerLocation {
				t.Errorf("Expected Location %s, got %s", tt.fields.headerLocation, res.Header.Get("Location"))
			}
			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.fields.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.fields.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
