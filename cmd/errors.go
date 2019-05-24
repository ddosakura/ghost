package cmd

import (
	"errors"
)

// errors
var (
	// --- Common ---
	//ErrUnknowSubCmd      = errors.New("Unknow sub-command")
	ErrUnknowServiceSign = errors.New("Unknow sign of sub-command `service`")
	ErrModelVersion      = errors.New("Unable to process this version of data")

	// -- Master ---
	// -- Node ---
	// -- Client ---
)
