package logger

import (
	"log"
	"os"
)

// NewLogger - создает новый логер
func NewLogger() *log.Logger {
	return log.New(os.Stdout, `gofermart | `, log.LstdFlags|log.Lshortfile)
}
