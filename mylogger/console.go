package mylogger

import (
	"fmt"
	"time"
)

//往终端写日志相关内容


type Logger struct {
	Level Loglevel
}


func Newconsolelog(levelstr string) Logger {
	level, err := ParseLogLevel(levelstr)
	if err != nil {
		panic(err)
	}
	return Logger{
		Level: level,
	}
}

func (l Logger) enable(loglevel Loglevel) bool {
	return loglevel>=l.Level
}

func (l Logger) Debug(format string,a...interface{}) {
	if l.enable(DEBUG) {
		msg:=fmt.Sprintf(format,a...)
		now := time.Now().Format("2006-01-02 15:04:05")
		funcName,fileName,lineNo:=getinfo(2)
		fmt.Printf("[%s] [DEBUG] [%s:%s:%d] %s\n", now,funcName,fileName,lineNo, msg)

	}

}

func (l Logger) Trace(format string,a...interface{}) {
	if l.enable(TRACE) {
		msg:=fmt.Sprintf(format,a...)
		now := time.Now().Format("2006-01-02 15:04:05")
		funcName,fileName,lineNo:=getinfo(2)
		fmt.Printf("[%s] [TRACE] [%s:%s:%d] %s\n", now,funcName,fileName,lineNo, msg)
	}
}

func (l Logger) Info(format string,a...interface{}) {
	if l.enable(INFO) {
		msg:=fmt.Sprintf(format,a...)
		now := time.Now().Format("2006-01-02 15:04:05")
		funcName,fileName,lineNo:=getinfo(2)
		fmt.Printf("[%s] [INFO] [%s:%s:%d] %s\n", now,funcName,fileName,lineNo, msg)
	}
}

func (l Logger) Warning(format string,a...interface{}) {
	if l.enable(WARNING) {
		msg:=fmt.Sprintf(format,a...)
		now := time.Now().Format("2006-01-02 15:04:05")
		funcName,fileName,lineNo:=getinfo(2)
		fmt.Printf("[%s] [WARNING] [%s:%s:%d] %s\n", now,funcName,fileName,lineNo, msg)
	}
}
func (l Logger) Error(format string,a...interface{}) {
	if l.enable(ERROR) {
		msg:=fmt.Sprintf(format,a...)
		now := time.Now().Format("2006-01-02 15:04:05")
		funcName,fileName,lineNo:=getinfo(2)
		fmt.Printf("[%s] [ERROR] [%s:%s:%d] %s\n", now,funcName,fileName,lineNo, msg)
	}
}
func (l Logger) Fatal(format string,a...interface{}) {
	if l.enable(FATAL) {
		msg:=fmt.Sprintf(format,a...)
		now := time.Now().Format("2006-01-02 15:04:05")
		funcName,fileName,lineNo:=getinfo(2)
		fmt.Printf("[%s] [FATAL] [%s:%s:%d] %s\n", now,funcName,fileName,lineNo, msg)
	}
}
