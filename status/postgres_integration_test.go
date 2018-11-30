// +build integration

package status

import (
	"reflect"
	"testing"
	"time"
)

func Test_New(t *testing.T) {
	tt := []struct {
		name    string
		connStr string
		wantErr bool
	}{
		{
			name:    "should pass if the connStr is valid",
			connStr: "postgres://postgres:postgres@localhost:5432/notfy_db?sslmode=disable",
			wantErr: false,
		},
		{
			name:    "should return error if the connStr is invalid",
			connStr: "postgres://postgres:postgres@localhost:9999/notfy",
			wantErr: true,
		},
	}
	for _, tst := range tt {
		t.Run(tst.name, func(t *testing.T) {
			_, err := NewPostgresStorage(tst.connStr)
			if tst.wantErr && err == nil {
				t.Errorf("New(): got no error, but expected an error")
			}
			if !tst.wantErr && err != nil {
				t.Errorf("New(): got error %v, but expected none", err)
			}
		})
	}
}

func Test_InsertAndGet(t *testing.T) {
	connStr := "postgres://postgres:postgres@localhost:5432/notfy_db?sslmode=disable"
	db, err := NewPostgresStorage(connStr)
	if err != nil {
		t.Errorf("cannot create storage: %v", err)
	}

	now := time.Now()
	info := MakeInfo(0, Queued)
	info.SetCreatedAt(now)
	info.SetLastUpdatedAt(now)
	ret, err := db.insert(info)
	if err != nil {
		t.Errorf("failed to insert: %v", err)
	}
	got, err := db.get(ret.ID())
	if err != nil {
		t.Errorf("failed to get: %v", err)
	}
	if !reflect.DeepEqual(got, ret) {
		t.Errorf("got %v, but inserted %v", got, ret)
	}
}
