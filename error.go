package client

// Error is used to represent errors returned by the API.
// It contains error message and HTTP code.
type Error interface {
	error
	Code() int
}

type resError struct {
	code int
	err  string
}

func (err *resError) Error() string { return err.err }
func (err *resError) Code() int     { return err.code }
