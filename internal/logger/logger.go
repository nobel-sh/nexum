package logger

import (
	"fmt"
	"os"
	"time"
)

type Logger struct {
	file *os.File
}

func (l *Logger) log(level, format string, args ...interface{}) {
	timestamp := time.Now().Format(time.RFC3339)
	message := fmt.Sprintf(format, args...)
	logEntry := fmt.Sprintf("[%s] %s: %s\n", timestamp, level, message)
	l.file.WriteString(logEntry)
	fmt.Print(logEntry)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log("INFO", format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log("ERROR", format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log("FATAL", format, args...)
	os.Exit(1)
}
