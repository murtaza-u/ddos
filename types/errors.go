package types

import "errors"

var newE = errors.New

var (
	ErrMarshaling        = newE("failed to marshal manifest")
	ErrUnmarshaling      = newE("failed to unmarshal manifest")
	ErrUnsupportedMethod = newE("method not supported or specified")
	ErrMethodNotAllowed  = newE("delete operation not allowed")
	ErrKeyNotFound       = newE("key not found")
	ErrKeyModified       = newE("the key has been modified after you read it")
	ErrGeneratingUID     = newE("failed to generate UID")
	ErrIdNotSet          = newE("id not set")
	ErrCanceled          = newE("context canceled")
)
