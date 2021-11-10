package mylogger

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"
)

/*
支持不同地方的输出
日志级别 debug trace  info warning error fatal
支持开关控制
包含时间，行号，文件名，日志级别，日志信息
日志文件切割
*/
type Loglevel uint16

//接口
type LoggerIO interface {
	Debug(format string, a ...interface{})
	Trace(format string, a ...interface{})
	Info(format string, a ...interface{})
	Warning(format string, a ...interface{})
	Error(format string, a ...interface{})
	Fatal(format string, a ...interface{})
}

const (
	UNKNOW Loglevel = iota
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

func ParseLogLevel(s string) (Loglevel, error) {
	s = strings.ToLower(s)
	switch s {
	case "debug":
		return DEBUG, nil
	case "trace":
		return TRACE, nil
	case "info":
		return INFO, nil
	case "waring":
		return WARNING, nil
	case "error":
		return ERROR, nil
	case "fatal":
		return FATAL, nil
	default:
		err := errors.New("无效的日志级别")
		return UNKNOW, err

	}
}

func getloglvstring(lv Loglevel) string {
	switch lv {
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	case  INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	}
	return "DEFAULT"
}

func getinfo(n int) (funcName, fileName string, lineNO int) {
	pc, file, lineNO, ok := runtime.Caller(n)
	if !ok {
		fmt.Println("runtime caller failed")
		return
	}
	funcName = runtime.FuncForPC(pc).Name()
	fileName = path.Base(file)
	return

}
