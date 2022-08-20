package tools


// Error401 - пользователь не авторизован
type Error401 struct {
	message string
}

func (e Error401) Error() string {
	return e.message
}

func NewError401(message string) *Error401 {
	return &Error401{
		message: message,
	}
}

