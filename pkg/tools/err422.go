package tools

// Error422 - не верный формат заказа
type Error422 struct {
	message string
}

func (e Error422) Error() string {
	return e.message
}

func NewError422(message string) *Error422 {
	return &Error422{
		message: message,
	}
}
