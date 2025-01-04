package calendar

import "errors"

var (
	ErrNotFound   = errors.New("calendar not found")
	ErrRepository = errors.New("repository error")
)
