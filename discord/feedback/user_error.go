package feedback

import "fmt"

type UserError struct {
	Message string
}

func (u UserError) Error() string {
	return u.Message
}

type WrappedUserError struct {
	UserError UserError
	Cause     error
}

func (w WrappedUserError) Error() string {
	return fmt.Sprintf("%s (%s)", w.UserError, w.Cause)
}

func Wrap(userError UserError, cause error) WrappedUserError {
	return WrappedUserError{
		UserError: userError,
		Cause:     cause,
	}
}
