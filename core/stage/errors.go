package stage

type ErrorCategory string

const (
	ErrIO             ErrorCategory = "ErrIO"
	ErrStorageCorrupt ErrorCategory = "ErrStorageCorrupt"
)
