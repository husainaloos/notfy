package email

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestPostgressStorageInsert(t *testing.T) {
	e, _ := New(0, "myself@myself.com", []string{"to@to.com"}, []string{"cc@cc.com", "cc2@cc.com"}, []string{}, "subject", "body")
	t1, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	t2, _ := time.Parse(time.RFC3339, "2007-02-02T15:04:05Z")
	e.AddStatusEvent(MakeStatusEvent(Queued, t1))
	e.AddStatusEvent(MakeStatusEvent(SentSuccessfully, t2))
	connStr := "postgres://postgres:postgres@localhost/notfy?sslmode=disable"
	pg, err := NewPostgresStorage(connStr)
	if err != nil {
		t.Fatalf("failed to create connection: %v", err)
	}
	got, err := pg.insert(context.Background(), e)
	if err != nil {
		t.Fatalf("failed to insert email: %v", err)
	}
	if got.ID() <= 0 {
		t.Fatalf("email id is %d, but expected it to be greater than 0", e.ID())
	}
	if !reflect.DeepEqual(got.From(), e.From()) {
		t.Fatalf("got %v, but expected %v", got.From(), e.From())
	}

	if !reflect.DeepEqual(got.To(), e.To()) {
		t.Fatalf("got %v, but expected %v", got.To(), e.To())
	}

	if !reflect.DeepEqual(got.CC(), e.CC()) {
		t.Fatalf("got %v, but expected %v", got.CC(), e.CC())
	}

	if !reflect.DeepEqual(got.BCC(), e.BCC()) {
		t.Fatalf("got %v, but expected %v", got.BCC(), e.BCC())
	}

	if !reflect.DeepEqual(got.Subject(), e.Subject()) {
		t.Fatalf("got %v, but expected %v", got.Subject(), e.Subject())
	}

	if !reflect.DeepEqual(got.Body(), e.Body()) {
		t.Fatalf("got %v, but expected %v", got.Body(), e.Body())
	}

	if !reflect.DeepEqual(got.StatusHistory(), e.StatusHistory()) {
		t.Fatalf("got %v, but expected %v", got.StatusHistory(), e.StatusHistory())
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
	expect, _ := New(1, "myself@myself.com", []string{"to@to.com"}, []string{"cc@cc.com", "cc2@cc.com"}, []string{}, "subject", "body")
	t1, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	t2, _ := time.Parse(time.RFC3339, "2007-02-02T15:04:05Z")
	expect.AddStatusEvent(MakeStatusEvent(Queued, t1))
	expect.AddStatusEvent(MakeStatusEvent(SentSuccessfully, t2))

	if !reflect.DeepEqual(e, expect) {
		t.Fatalf("got %v, but expected %v", e, expect)
	}
}

func TestPostgressStorageUpdate(t *testing.T) {
	connStr := "postgres://postgres:postgres@localhost/notfy?sslmode=disable"
	pg, err := NewPostgresStorage(connStr)
	if err != nil {
		t.Fatalf("failed to create connection: %v", err)
	}
	e, _ := New(2, "myself@myself.com", []string{"to@to.com"}, []string{"cc@cc.com", "cc2@cc.com"}, []string{}, "subject", "body")
	now := time.Now()
	e.AddStatusEvent(MakeStatusEvent(Queued, now))
	e.AddStatusEvent(MakeStatusEvent(SentSuccessfully, now))

	got, ok, err := pg.update(context.Background(), e)
	if err != nil {
		t.Fatalf("failed to update email: %v", err)
	}
	if !ok {
		t.Fatal("email does not exists, but the email should exists")
	}
	if !reflect.DeepEqual(got, e) {
		t.Fatalf("got %v, but expected %v", got, e)
	}

	got, ok, err = pg.get(context.Background(), 2)
	if err != nil || !ok {
		t.Fatalf("ok=%t, err=%v", ok, err)
	}

	if !reflect.DeepEqual(got, e) {
		t.Fatalf("got %v, but expected %v", got, e)
	}
}
