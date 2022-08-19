package tools


// Error402 - не достаточно средств
type Error402 struct {
	message string
}

func (e Error402) Error() string {
	return e.message
}

func NewError402(message string) error {
	return &Error402{
		message: message,
	}
}

