package review

import (
	"errors"
	"fmt"
)

var (
	ErrReview               = errors.New("review error")
	ErrNotFound             = fmt.Errorf("%w: review not found", ErrReview)
	ErrRepository           = fmt.Errorf("%w: repository error", ErrReview)
	ErrOrderNotCompleted    = fmt.Errorf("%w: order not completed", ErrReview)
	ErrOrderAlreadyReviewed = fmt.Errorf("%w: order already reviewed", ErrReview)
)
