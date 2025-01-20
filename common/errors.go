package common

import "errors"

var (
	ErrorNoItem = errors.New("items must have at least 1 item!")
)
