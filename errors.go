package rollbar

import (
	"fmt"
)

// ErrHttpError is an HTTP error status code as defined by
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html
type ErrHttpError int

// Error implements the error interface.
func (e ErrHttpError) Error() string {
	return fmt.Sprintf("rollbar: service returned status: %d", e)
}
