package cli

import "errors"

var (
	// ErrExitSuccess represents an error that can be returned from an argument Flag.OnSetFunc func
	// if returned, program exits normally, error return 0
	ErrExitSuccess = errors.New("success")
)
