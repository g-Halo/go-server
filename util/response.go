package util

type Response struct {
	Code int
	Data interface{}
	Msg  string
}

const (
	Success = 1
	Fail    = -1
)
