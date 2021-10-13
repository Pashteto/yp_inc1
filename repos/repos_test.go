package repos

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Test_repository_Get(t *testing.T) {
	type fields struct {
		//		Client redis.Cmdable
		connPool *pgxpool.Pool
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{
				connPool: tt.fields.connPool,
			}
			r.Ping(context.Background())
			/*	got, err := r.Get(tt.args.ctx, tt.args.key)
				if (err != nil) != tt.wantErr {
					t.Errorf("repository.Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("repository.Get() = %v, want %v", got, tt.want)
				}*/
		})
	}
}

func Test_repository_Set(t *testing.T) {
	type fields struct {
		connPool *pgxpool.Pool
	}
	type args struct {
		ctx   context.Context
		key   string
		value interface{}
		exp   time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{
				connPool: tt.fields.connPool,
			}
			r.Ping(context.Background())
			/*
				if err := r.Set(tt.args.ctx, tt.args.key, tt.args.value, tt.args.exp); (err != nil) != tt.wantErr {
					t.Errorf("repository.Set() error = %v, wantErr %v", err, tt.wantErr)
				}*/
		})
	}
}
