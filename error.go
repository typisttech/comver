package comver

import "fmt"

type stringError string

func (s stringError) Error() string {
	return string(s)
}

type ParseError struct {
	original string
	wrapped  error
}

func (e ParseError) Error() string {
	return fmt.Sprintf("error parsing version string %q", e.original)
}

func (e ParseError) Unwrap() error {
	return e.wrapped
}

func (e ParseError) Original() string {
	return e.original
}
