package log

import (
	"errors"
	"log"
	"os"
	"strconv"
)

type LogLevel int

const (
	NONE  LogLevel = 0
	ERROR LogLevel = 1
	WARN  LogLevel = 2
	INFO  LogLevel = 3
	DEBUG LogLevel = 4
)

func LogLevelFromStr(str string) (LogLevel, error) {
	int, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		return 0, err
	}

	if int < 0 || int > 4 {
		return 0, errors.New("Invalid log level")
	}

	return LogLevel(int), nil
}

type Logger struct {
	log   *log.Logger
	level LogLevel
}

func NewLogger(level LogLevel) Logger {
	logger := Logger{
		log:   log.New(os.Stdout, "", log.Ldate|log.LUTC),
		level: level,
	}

	return logger
}

func (logger *Logger) Debug(args ...any) {
	if logger.level >= DEBUG {
		logger.log.SetPrefix("[DEBUG] ")
		logger.log.SetFlags(log.Ldate | log.LUTC | log.Lshortfile)
		logger.log.Print(args...)
		logger.log.SetFlags(log.Ldate | log.LUTC)
	}
}

func (logger *Logger) Info(args ...any) {
	if logger.level >= INFO {
		logger.log.SetPrefix("[INFO] ")
		logger.log.Print(args...)
	}
}

func (logger *Logger) Warn(args ...any) {
	if logger.level >= WARN {
		logger.log.SetPrefix("[WARN] ")
		logger.log.Print(args...)
	}
}

func (logger *Logger) Error(args ...any) {
	if logger.level >= ERROR {
		logger.log.SetPrefix("[ERROR] ")
		logger.log.Print(args...)
	}
}
