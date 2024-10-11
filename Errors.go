package searchquerylexer

import "errors"

var (
	ErrInvalidEscapeSequence error = errors.New("invalid escape sequence")
	ErrInvalidConnective     error = errors.New("invalid connective")

	ErrInvalidConfigComparator error = errors.New("invalid comparator config")
	ErrInvalidConfigConnective error = errors.New("invalid connective config")
)
