package tools

// 错误处理

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Msg  error
}

func (s StatusError) Error() string {
	return s.Msg.Error()
}

func (s StatusError) Status() int {
	return s.Code
}
