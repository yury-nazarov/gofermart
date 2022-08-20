package tools

// Error400 - не верный формат запроса
type Error400 struct {
	message string
}

func (e Error400) Error() string {
	return e.message
}

func NewError400(message string) *Error400 {
	return &Error400{
		message: message,
	}
}

