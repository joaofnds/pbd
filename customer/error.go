package customer

import "errors"

var (
	ErrNotFound   = errors.New("customer not found")
	ErrRepository = errors.New("repository error")
)
