package errors

type ErrAnswer struct {
	mess string
}

func NewErrAnswer(mess string) *ErrAnswer {
	return &ErrAnswer{mess: mess}
}

func (err *ErrAnswer) Error() string {
	return err.mess
}
