package booking

import (
	"errors"
	"fmt"
)

var (
	ErrBooking = errors.New("booking error")

	ErrRepository = fmt.Errorf("%w: repository error", ErrBooking)
	ErrNotFound   = fmt.Errorf("%w: not found", ErrRepository)

	// field validation errors
	ErrMissingWorkerID   = fmt.Errorf("%w: missing worker id", ErrBooking)
	ErrMissingCalendarID = fmt.Errorf("%w: missing calendar id", ErrBooking)
	ErrMissingEventID    = fmt.Errorf("%w: missing event id", ErrBooking)
	ErrMissingStartsAt   = fmt.Errorf("%w: missing starts at", ErrBooking)
	ErrMissingEndsAt     = fmt.Errorf("%w: missing ends at", ErrBooking)
	ErrStartAfterEnd     = fmt.Errorf("%w: starts at must be before ends at", ErrBooking)
	ErrStartEqualEnd     = fmt.Errorf("%w: starts at must be before ends at", ErrBooking)
	ErrStartAfterNow     = fmt.Errorf("%w: starts at must be in the future", ErrBooking)
	ErrInvalidDuration   = fmt.Errorf("%w: invalid duration", ErrBooking)

	// booking errors
	ErrEventNotAvailable  = fmt.Errorf("%w: event is not available", ErrBooking)
	ErrNotWithinEventTime = fmt.Errorf("%w: request is not within event time", ErrBooking)
)
