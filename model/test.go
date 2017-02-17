package model

import (
	"time"
)

// A Test represents a single run of a test.
type Test struct {
	Name            string
	FixtureName     string
	Result          string
	TimeStarted     time.Time
	TimeEnded       time.Time
	DurationSeconds float64
	ErrorMessage    *string
	StackTrace      *string
}
