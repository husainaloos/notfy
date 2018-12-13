package email

import (
	"context"
	"testing"
	"time"
)

func TestPostgressStorageInsert(t *testing.T) {
	e, _ := New(0, "myself@myself.com", []string{"to@to.com"}, []string{"cc@cc.com", "cc2@cc.com"}, []string{}, "subject", "body")
	e.AddStatusEvent(MakeStatusEvent(Created, time.Now()))
	e.AddStatusEvent(MakeStatusEvent(Queued, time.Now()))
	connStr := "postgres://postgres:postgres@localhost/notfy?sslmode=disable"
	pg, err := NewPostgresStorage(connStr)
	if err != nil {
		t.Fatalf("failed to create connection: %v", err)
	}
	e, err = pg.insert(context.Background(), e)
	if err != nil {
		t.Fatalf("failed to insert email: %v", err)
	}
	if e.ID() <= 0 {
		t.Fatalf("email id is %d, but expected it to be greater than 0", e.ID())
	}
}

func TestPostgressStorageGet(t *testing.T) {
	connStr := "postgres://postgres:postgres@localhost/notfy?sslmode=disable"
	pg, err := NewPostgresStorage(connStr)
	if err != nil {
		t.Fatalf("failed to create connection: %v", err)
	}
	e, ok, err := pg.get(context.Background(), 1)
	if err != nil {
		t.Fatalf("failed to insert email: %v", err)
	}
	if !ok {
		t.Fatal("email does not exists, but the email should exists")
	}
	t.Fatal(e)
}
