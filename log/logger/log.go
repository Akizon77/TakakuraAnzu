package logger

import (
	"fmt"
	logs "github.com/Akizon77/TakakuraAnzu/log"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type Logger struct {
	logger *log.Logger
	file   *os.File
}

func (l Logger) Debug(v ...interface{}) {
	//fmt.Println(v)
}

func (l Logger) Info(v ...interface{}) {
	logs.Info(fmt.Sprint(v))

}

func (l Logger) Warn(v ...interface{}) {
	logs.Warn(fmt.Sprint(v))
}

func (l Logger) Error(v ...interface{}) {
	logs.Error(fmt.Sprint(v), nil)
}

func (l Logger) Debugf(format string, v ...interface{}) {
	//fmt.Println(v)
}

func (l Logger) Infof(format string, v ...interface{}) {
	//fmt.Println(v)
}

func (l Logger) Warnf(format string, v ...interface{}) {
	//fmt.Println(v)
}

func (l Logger) Errorf(format string, v ...interface{}) {
	//fmt.Println(v)
}

func (l Logger) Sync() error {
	return nil
}

func NewLogger() *Logger {
	return &Logger{}
}
func output(v ...interface{}) string {
	_, file, line, _ := runtime.Caller(3)
	files := strings.Split(file, "/")
	file = files[len(files)-1]

	logFormat := "%s %s:%d " + fmt.Sprint(v...) + "\n"
	date := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf(logFormat, date, file, line)
}
