package email

import "time"

type Status int

//go:generate stringer -type Status
const (
	Created Status = iota
	Queued
	AttemptedToSend
	SentSuccessfully
	FailedToSend
	QueuedForRetry
	Dead
)

type StatusEvent struct {
	status Status
	at     time.Time
}

func NewStatusEvent(status Status, at time.Time) StatusEvent {
	return StatusEvent{status, at}
}

type StatusHistory []StatusEvent
