package atm

import "errors"

//
var (
	ErrEmpty             = errors.New("")
	ErrCommandFormat     = errors.New("Format: Command format invalid")
	ErrUsernameIncorrect = errors.New("Format: Username incorrect")
)
