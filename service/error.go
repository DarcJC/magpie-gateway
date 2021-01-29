package service

type Error struct {
	s string
}

func (e *Error) Error() string {
	return e.s
}

func NewError(desc string) error {
	return &Error{
		s: desc,
	}
}
