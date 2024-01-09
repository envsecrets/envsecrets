package commons

import "errors"

const (
	NONCE_LEN = 24
	KEY_BYTES = 32
)

var (
	ErrNoServerKey = errors.New("SERVER_SYMMETRIC_KEY is not set")
)
