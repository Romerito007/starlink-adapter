package codes

type Code int32

const (
	OK               Code = 0
	DeadlineExceeded Code = 4
	Unimplemented    Code = 12
	Unavailable      Code = 14
)
