package platform

import "errors"

var (
	ErrUnsupported = errors.New("unsupported")
	ErrUnknownDst  = errors.New("unknown destination")
)
