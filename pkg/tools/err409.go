package tools

// Error409 - конфликт
type Error409 struct {
	message string
}

func (e Error409) Error() string {
	return e.message
}

func NewError409(message string) error {
	return &Error409{
		message: message,
	}
}

