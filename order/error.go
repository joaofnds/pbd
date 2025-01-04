package order

import (
	"errors"
	"fmt"
)

var (
	ErrOrder = errors.New("order error")

	ErrRepository = fmt.Errorf("%w: repository error", ErrOrder)
	ErrNotFound   = fmt.Errorf("%w: not found", ErrRepository)

	// order errors
	ErrEventNotAvailable  = fmt.Errorf("%w: event is not available", ErrOrder)
	ErrNotWithinEventTime = fmt.Errorf("%w: request is not within event time", ErrOrder)
)
