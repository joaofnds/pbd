package event

import "errors"

var (
	ErrNotFound      = errors.New("event not found")
	ErrRepository    = errors.New("repository error")
	ErrTimeSlotTaken = errors.New("time slot is taken")
)
