package status

import "time"

//go:generate stringer -type=SendStatus

// SendStatus is the status of sent notifications
type SendStatus int

const (
	// Sent is any notification that has been fulfilled
	Sent SendStatus = iota
	// Failed is any notification that failed to be sent
	Failed
	// Queued is any queued notifications that are in the process of being sent
	Queued
)

// Info is the details of the status of a notification
type Info struct {
	id           int
	status       SendStatus
	createdAt    time.Time
	lastUpdateAt time.Time
}

// ID is the id of the status
func (s Info) ID() int { return s.id }

// Status is the status
func (s Info) Status() SendStatus { return s.status }

// CreatedAt is when the status has been created
func (s Info) CreatedAt() time.Time { return s.createdAt }

// LastUpdateAt is when the status has been updated last
func (s Info) LastUpdateAt() time.Time { return s.lastUpdateAt }

// MakeInfo creates a new Info
func MakeInfo(id int, status SendStatus, createdAt, lastUpdateAt time.Time) Info {
	return Info{
		id:           id,
		status:       status,
		createdAt:    createdAt.UTC(),
		lastUpdateAt: lastUpdateAt.UTC(),
	}
}
