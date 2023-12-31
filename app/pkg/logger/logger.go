package logger

import (
	"boost-my-skills-bot/pkg/utils"
	"fmt"
)

const (
	Red   = "\033[31m"
	Green = "\033[32m"
	Reset = "\033[0m"
)

type Logger struct {
}

func InitLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Infof(message string, args ...any) {
	formatedMessage := fmt.Sprintf(message, args...)
	fmt.Printf("%s %s [INFO] %s %s\n", utils.GetFormatedTime(), Green, Reset, formatedMessage)
}

func (l *Logger) Errorf(message string, args ...any) {
	formatedMessage := fmt.Sprintf(message, args...)
	fmt.Printf("%s %s [ERROR] %s %s\n", utils.GetFormatedTime(), Red, Reset, formatedMessage)
}
