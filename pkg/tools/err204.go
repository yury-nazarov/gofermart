package tools

// Error204 нет данных
type Error204 struct {
	message string
}

func (e Error204) Error() string {
	return e.message
}

func NewError204(message string) *Error204 {
	return &Error204{
		message: message,
	}
}