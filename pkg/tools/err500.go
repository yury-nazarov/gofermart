package tools


// Error500 - ошибка сервера
type Error500 struct {
	message string
}

func (e Error500) Error() string {
	return e.message
}

func NewError500(message string) *Error500 {
	return &Error500{
		message: message,
	}
}
