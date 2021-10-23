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
		})
	}
}
