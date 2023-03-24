package log

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarning
	LevelError
)

var (
	logInfoPath  = ""
	logWarnPath  = ""
	logErrorPath = ""
)

func init() {
	path := GetPath()
	if len(path) == 0 {
		return
	}
	logInfoPath = path + "\\base_logs\\info.log"
	logWarnPath = path + "\\base_logs\\warn.log"
	logErrorPath = path + "\\base_logs\\error.log"
	GetLogger()
}

type Logger struct {
	infoChan  chan string
	warnChan  chan string
	errorChan chan string
	infoFile  *os.File
	warnFile  *os.File
	errorFile *os.File
	once      sync.Once
}

var LoggerSingle *Logger
var ctx, cancel = context.WithCancel(context.Background())

func GetLogger() *Logger {
	if LoggerSingle == nil {
		LoggerSingle = &Logger{}
		LoggerSingle.init()
	}
	return LoggerSingle
}
func GetPath() string {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}
	return path
}
func (l *Logger) init() {
	l.once.Do(func() {
		l.infoChan = make(chan string, 1)
		l.warnChan = make(chan string, 1)
		l.errorChan = make(chan string, 1)
		infoFile, err := os.OpenFile(logInfoPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("failed to create info.log: %v", err))
		}
		l.infoFile = infoFile
		warnFile, err := os.OpenFile(logWarnPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("failed to create warn.log: %v", err))
		}
		l.warnFile = warnFile
		errorFile, err := os.OpenFile(logErrorPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("failed to create error.log: %v", err))
		}
		l.errorFile = errorFile
		go l.writeToLogFile(ctx, l.infoChan, l.infoFile, "INFO: ")
		go l.writeToLogFile(ctx, l.warnChan, l.warnFile, "WARN: ")
		go l.writeToLogFile(ctx, l.errorChan, l.errorFile, "ERROR: ")
	})
}
func (l *Logger) StopLogger() {
	cancel()
	close(l.warnChan)
	close(l.infoChan)
	close(l.errorChan)
	l.infoFile.Close()
	l.errorFile.Close()
	l.warnFile.Close()
}
func (l *Logger) writeToLogFile(ctx context.Context, ch chan string, file *os.File, prefix string) {
	for {
		select {
		case message := <-ch:
			now := time.Now().Format("2006-01-02 15:04:05")
			logMessage := fmt.Sprintf("%s %s%s\n", now, prefix, message)
			_, err := file.WriteString(logMessage)
			if err != nil {
				log.Printf("failed to write log message: %v", err)
			}
		case <-ctx.Done():
			fmt.Println(prefix, "close")
			return
		}
	}
}

func (l *Logger) log(level LogLevel, message string) {
	switch level {
	case LevelDebug:
		log.Printf("[DEBUG] %s", message)
	case LevelInfo:
		l.infoChan <- message
	case LevelWarning:
		l.warnChan <- message
	case LevelError:
		l.errorChan <- message
	}
}

func (l *Logger) Debug(message string) {
	l.log(LevelDebug, message)
}

func (l *Logger) Info(message string) {
	l.log(LevelInfo, message)
}

func (l *Logger) Warning(message string) {
	l.log(LevelWarning, message)
}

func (l *Logger) Error(message string) {
	l.log(LevelError, message)
}
