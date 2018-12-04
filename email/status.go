package email

import "time"

type Status int

//go:generate stringer -type Status
const (
	Created Status = iota
	Queued
	SentSuccessfully
	FailedAttemptToSend
	QueuedForRetry
	Dead
)

type StatusEvent struct {
	status Status
	at     time.Time
}

func MakeStatusEvent(status Status, at time.Time) StatusEvent {
	return StatusEvent{status, at.UTC()}
}

func (se StatusEvent) Status() Status { return se.status }
func (se StatusEvent) At() time.Time  { return se.at }

type StatusHistory []StatusEvent
