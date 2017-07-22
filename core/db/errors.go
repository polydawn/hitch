package db

type ErrorCategory string

const (
	ErrNotFound       ErrorCategory = "ErrNotFound"
	ErrIO             ErrorCategory = "ErrIO"
	ErrStorageCorrupt ErrorCategory = "ErrStorageCorrupt"
)
