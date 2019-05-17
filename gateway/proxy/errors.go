package gateproxy

import (
	"errors"
)

var (
	// Common

	// HTTP
	errAuthMethodUnsupported = errors.New("auth method unsupported")
	errAuthNull              = errors.New("user/pass is null")

	// SOCKET
)
