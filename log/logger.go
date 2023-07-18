package log

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type LogLevel uint

type Logger struct {
	*log.Logger
	name string
	addr string
}

const (
	DEBUG LogLevel = iota
	INFO
	ERROR
	DISABLE
)

var (
	debugLogger = new("\033[33m[DEBUG]\033[0m ")
	infoLogger  = new("\033[34m[INFO]\033[0m ")
	errorLogger = new("\033[31m[ERROR]\033[0m ")
	Debugf      = debugLogger.Printf
	Debug       = debugLogger.Println
	Infof       = infoLogger.Printf
	Info        = infoLogger.Println
	Errorf      = errorLogger.Printf
	Error       = errorLogger.Println

	loggers = []*Logger{errorLogger, infoLogger, debugLogger}
	mu      sync.Mutex
)

func new(predix string) *Logger {
	return &Logger{Logger: log.New(os.Stdout, predix, log.LstdFlags|log.Lshortfile)}
}

func WithServer(name string, addr string) {
	for _, logger := range loggers {
		logger.name = name
		logger.addr = addr
	}
}

func (logger *Logger) Printf(format string, args ...interface{}) {
	if logger.name != "" {
		format = "[%s(%s)]" + format
		args = append([]interface{}{logger.name, logger.addr}, args...)
	}
	logger.Logger.Printf(format, args...)
}
func (logger *Logger) Println(msg string) {
	if logger.name != "" {
		msg = fmt.Sprintf("[ %s (%s) ]", logger.name, logger.addr) + msg
	}
	logger.Logger.Println(msg)
}

func SetLevel(lv LogLevel) {
	mu.Lock()
	defer mu.Unlock()
	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}
	if lv > ERROR {
		errorLogger.SetOutput(ioutil.Discard)
	}

	if lv > INFO {
		infoLogger.SetOutput(ioutil.Discard)
	}

	if lv > DEBUG {
		debugLogger.SetOutput(ioutil.Discard)
	}
}
