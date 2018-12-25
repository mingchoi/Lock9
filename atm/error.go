package atm

import "errors"

//
var (
	ErrCommandFormat     = errors.New("Format: Command format invalid")
	ErrUsernameIncorrect = errors.New("Format: Username incorrect")
)
