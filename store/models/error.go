package models

type ModelError struct {
    s string
}

func (e *ModelError) Error() string {
    return e.s
}

func NewError(desc string) error {
    return &ModelError{
        s: desc,
    }
}

