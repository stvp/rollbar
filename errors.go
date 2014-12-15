package rollbar

import (
	"fmt"
)

// ErrNoToken is returned when trying to send a message to the rollbar service,
// while not having the token initialized.
type ErrNoToken struct{}

// Error implements the error interface.
func (e ErrNoToken) Error() string {
	return "rollbar: No token set."
}

type ErrHttpError int

// Error implements the error interface.
func (e ErrHttpError) Error() string {
	return fmt.Sprintf("rollbar: The rollbar service replied with %d http status code.", e)
}
