package sh

import (
	"errors"
)

// errors
var (
	// init
	ErrDirIsFile = errors.New("The necessary ghost system directory is a file")

	// vm
	ErrUnknowExpr = errors.New("Unknow Expr")
)
