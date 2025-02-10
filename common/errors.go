package common

import "errors"

var (
	ErrorNoItem  = errors.New("items must have at least 1 item!")
	ErrorNoStock = errors.New("some item is not in stock")
)
