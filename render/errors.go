package render

import "errors"

var (
	ErrNoEngine = errors.New("invalid engine specified")
	ErrNoPage   = errors.New("no page specified")
	ErrExecPage = errors.New("can not execute template")
)
