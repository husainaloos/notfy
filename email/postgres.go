package email

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type pgStatusEvent struct {
	Status int32     `json:"status"`
	At     time.Time `json:"at"`
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{db}, nil
}

func (s *PostgresStorage) insert(ctx context.Context, e Email) (Email, error) {
	statusEvents := []pgStatusEvent{}
	for _, v := range e.StatusHistory() {
		statusEvents = append(statusEvents, pgStatusEvent{int32(v.Status()), v.At()})
	}
	bin, err := json.Marshal(&statusEvents)
	if err != nil {
		return Email{}, fmt.Errorf("cannot json.Marshal status: %v", err)
	}
	emailID := 0
	query := `INSERT INTO notfy.email ("from", "to", cc, bcc, subject, body, status_events) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING email_id`
	err = s.db.QueryRowContext(ctx, query, e.StringFrom(), pq.Array(e.StringTo()), pq.Array(e.StringCC()), pq.Array(e.StringBCC()), e.Subject(), e.Body(), string(bin)).Scan(&emailID)
	if err != nil {
		return Email{}, err
	}
	e.SetID(emailID)
	return e, nil
}

func (s *PostgresStorage) get(ctx context.Context, id int) (Email, bool, error) {
	query := `SELECT (email_id, "from", "to", cc, bcc, subject, body, status_events) FROM notfy.email WHERE email_id = $1`
	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return Email{}, true, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var from, subject, body, dbStatusEvent string
		var pqTo, pqCC, pqBCC []sql.NullString
		if err := rows.Scan(&id, &from, pq.Array(&pqTo), pq.Array(&pqCC), pq.Array(&pqBCC), &subject, &body, &dbStatusEvent); err != nil {
			c, _ := rows.Columns()
			return Email{}, true, fmt.Errorf("cannot scan row: %v. columns: %v", err, c)
		}

		var to, cc, bcc []string
		for _, r := range pqTo {
			var str string
			r.Scan(&str)
			to = append(to, str)
		}
		for _, r := range pqCC {
			var str string
			r.Scan(&str)
			cc = append(cc, str)
		}
		for _, r := range pqBCC {
			var str string
			r.Scan(&str)
			bcc = append(bcc, str)
		}
		se := []pgStatusEvent{}
		if err := json.Unmarshal([]byte(dbStatusEvent), &se); err != nil {
			return Email{}, true, fmt.Errorf("cannot json.Unmarshal status_events: %v and values is %s and from %s", err, dbStatusEvent, from)
		}
		e, err := New(id, from, to, cc, bcc, subject, body)
		if err != nil {
			return Email{}, true, fmt.Errorf("cannot build email: %v", err)
		}
		for _, v := range se {
			e.AddStatusEvent(MakeStatusEvent(Status(v.Status), v.At))
		}
		return e, true, nil
	}
	return Email{}, false, nil
}
