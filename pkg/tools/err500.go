package tools


// Error500 - ошибка сервера
type Error500 struct {
	message string
}

func (e Error500) Error() string {
	return e.message
}

func NewError500(message string) error {
	return &Error500{
		message: message,
	}
}
