package worker

import "errors"

var (
	ErrNotFound   = errors.New("worker not found")
	ErrRepository = errors.New("repository error")
)
