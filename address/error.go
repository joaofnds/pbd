package address

import (
	"errors"
	"fmt"
)

var (
	ErrAddress    = errors.New("address error")
	ErrNotFound   = fmt.Errorf("%w: %s", ErrAddress, "address not found")
	ErrRepository = fmt.Errorf("%w: %s", ErrAddress, "repository error")
)
