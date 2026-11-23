package lthash

import (
	"codeberg.org/reiver/go-erorr"
)

const (
	ErrBadStateSize = erorr.Error("lthash: bad state size — must be 2048 bytes")
)
