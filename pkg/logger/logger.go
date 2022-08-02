package logger

import (
	"fmt"
	"log"
	"os"
)

// NewLogger - создает новый логер
func NewLogger(prefix string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf(`%s | `, prefix), log.LstdFlags|log.Lshortfile)
}
