package cmd

import "time"

// LevelType provided to select log level
type LevelType string

const (
	NormalType LevelType = "noraml"
	DebugType  LevelType = "debug"
	InfoType   LevelType = "info" // default
	WarnType   LevelType = "warn"
	ErrorType  LevelType = "error"
	FatalType  LevelType = "fatal"
	PanicType  LevelType = "panic"
)

type ReasonType string

const (
	RTCreating          ReasonType = "creating"
	RTCreated           ReasonType = "created"
	RTUpdating          ReasonType = "updating" // default
	RTUpdated           ReasonType = "updated"
	RTDeleting          ReasonType = "deleting"
	RTDeleted           ReasonType = "deleted"
	RTPolicyCheck       ReasonType = "policy-check"
	RTPolicyCheckSucced ReasonType = "policy-success"
	RTPolicyCheckFailed ReasonType = "policy-failure"
)

type EventSummary struct {
	// FirstSeen time this event was first seen, aka 1 minute ago
	FirstSeen int

	// LastSeen time this event was last seen
	LastSeen int

	// Count number of times this event has occurred
	Count int

	//Source of event
	From string

	// SubObjectPath where in this object the event occurred
	SubObjectPath string

	// Type event error level, normal, warning, etc
	Type string

	// Reason what triggered this event, created, schedule, etc
	Reason string

	// Message help humans
	Message string
}

type EventDetail struct {
	// Time of event
	Time time.Time

	Type string

	//Source of event
	From string

	// SubObjectPath where in this object the event occurred
	SubObjectPath string

	Reason string

	Massage string
}
