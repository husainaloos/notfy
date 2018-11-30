package status

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db is not pingable: %v", err)
	}
	return &PostgresStorage{db}, nil
}

func (p *PostgresStorage) insert(info Info) (Info, error) {
	result, err := p.db.Exec("insert into notfy.status (code, created_at, last_update_at) values (?, ?, ?)",
		info.Status().String(), info.CreatedAt(), info.LastUpdateAt())
	if err != nil {
		return Info{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Info{}, fmt.Errorf("failed to get last inserted id: %v", err)
	}
	i := MakeInfo(int(id), info.Status())
	i.SetCreatedAt(info.CreatedAt())
	i.SetLastUpdatedAt(info.LastUpdateAt())
	return i, nil
}

func (p *PostgresStorage) update(info Info) (Info, error) {
	_, err := p.db.Exec("update notfy.status set code=?, created_at=?, last_update_at=? where status_id = ?",
		info.Status().String(), info.CreatedAt(), info.LastUpdateAt(), info.ID())
	if err != nil {
		return Info{}, fmt.Errorf("failed to update: %v", err)
	}
	return info, nil

}

func (p *PostgresStorage) get(id int) (Info, error) {
	rows, err := p.db.Query("select status_id, code, created_at, last_update_at from status where status_id = ?", id)
	if err != nil {
		return Info{}, fmt.Errorf("cannot get from db: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id           int
			code         int
			createdAt    time.Time
			lastUpdateAt time.Time
		)
		if err := rows.Scan(&id, &code, &createdAt, &lastUpdateAt); err != nil {
			return Info{}, fmt.Errorf("cannot read result: %v", err)
		}
		i := MakeInfo(id, SendStatus(code))
		i.SetCreatedAt(createdAt)
		i.SetLastUpdatedAt(lastUpdateAt)
		return i, nil
	}

	return Info{}, errStorageNotFound
}

func (p *PostgresStorage) Close() error {
	return p.db.Close()
}
