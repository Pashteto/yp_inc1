package filedb

import (
	"testing"

	"github.com/Pashteto/yp_inc1/config"
	"github.com/Pashteto/yp_inc1/repos"
)

func TestWriteAll(t *testing.T) {
	type args struct {
		rdb       repos.SetterGetter
		cfg       config.Config
		UsersList *[]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteAll(tt.args.rdb, tt.args.cfg, tt.args.UsersList); (err != nil) != tt.wantErr {
				t.Errorf("WriteAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
