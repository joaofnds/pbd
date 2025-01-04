package phone

import (
	"errors"
	"fmt"
)

var (
	ErrPhone      = errors.New("phone error")
	ErrNotFound   = fmt.Errorf("%w: %s", ErrPhone, "phone not found")
	ErrRepository = fmt.Errorf("%w: %s", ErrPhone, "repository error")
)
