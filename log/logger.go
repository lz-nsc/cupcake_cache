package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type LogLevel uint

const (
	DEBUG LogLevel = iota
	INFO
	ERROR
	DISABLE
)

var (
	debugLogger = log.New(os.Stdout, "\033[33m[DEBUG]\033[0m ", log.LstdFlags|log.Lshortfile)
	infoLogger  = log.New(os.Stdout, "\033[34m[INFO]\033[0m ", log.LstdFlags|log.Lshortfile)
	errorLogger = log.New(os.Stdout, "\033[31m[ERROR]\033[0m ", log.LstdFlags|log.Lshortfile)
	Debugf      = debugLogger.Printf
	Debug       = debugLogger.Println
	Infof       = infoLogger.Printf
	Info        = infoLogger.Println
	Errorf      = errorLogger.Printf
	Error       = errorLogger.Println

	loggers = []*log.Logger{errorLogger, infoLogger, debugLogger}
	mu      sync.Mutex
)

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
