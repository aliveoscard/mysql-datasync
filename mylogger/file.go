package mylogger

import (
	"fmt"
	"os"
	"path"
	"time"
)

//往文件里写日志

type FileLog struct {
	level      Loglevel
	filePath   string //日志文件路径
	fileName   string //日志文的文件名
	fileobj    *os.File
	errfileobj *os.File
	maxSize    int64
	logchan    chan *logmsg
}


type logmsg struct{
	level Loglevel
	msg string
	funcname string
	filename string
	line int
	timestamp string	
}


func NewFileLog(levelstr, fp, fn string, maxSize int64) *FileLog {
	level, err := ParseLogLevel(levelstr)
	if err != nil {
		panic(err)
	}
	f1 := &FileLog{
		level:    level,
		filePath: fp,
		fileName: fn,
		maxSize:  maxSize,
		logchan: make(chan *logmsg,50000),
	}
	err=f1.initfile()
	if err!=nil{
		panic(err)
	}
	return f1
}


//打开日志文件
func (f *FileLog) initfile() error {
	fullpath := path.Join(f.filePath, f.fileName)
	fileObj, err := os.OpenFile(fullpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("打开日志文件出错，err:", err)
		return err
	}
	fullpath= path.Join(f.filePath, "ERR"+f.fileName)
	errorObj, err := os.OpenFile(fullpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("打开日志文件出错，err:", err)
		return err
	}
	f.fileobj=fileObj
	f.errfileobj=errorObj
	// for i := 0; i < 5; i++ {
	// 	go f.WriteLogbg()  //后台写日志
	// }
	go f.WriteLogbg() 
	return nil
}

//关闭日志文件
func (f *FileLog)Close()  {
	f.fileobj.Close()
	f.errfileobj.Close()
}


//日志切割
func (f *FileLog)checksize(file *os.File) bool {
	fileinfo,err:=file.Stat()
	if err!=nil{
		fmt.Println("获取文件信息失败")
		return false
	}
	return fileinfo.Size()>=f.maxSize
}
func (f *FileLog)splitelog(file *os.File) (*os.File,error) {
	nowstr:=time.Now().Format("20060102150405000")
	fileinfo,_:=file.Stat()
	logname:=path.Join(f.filePath,fileinfo.Name())
	file.Close()
	newlogname:=fmt.Sprintf("%s.bak%s",logname,nowstr)
	os.Rename(logname,newlogname)
	fileObj, err := os.OpenFile(logname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("打开日志文件出错，err:", err)
		return nil,err
	}
	return fileObj,nil
}


func (f *FileLog) enable(loglevel Loglevel) bool {
	return loglevel >= f.level
}

//后台写日志
func (f *FileLog)WriteLogbg()  {
	for {
		if f.checksize(f.fileobj){ //切割日志文件
			newfile,_:=f.splitelog(f.fileobj)
			f.fileobj=newfile
			}

		
		select{
		case logtmp:=<-f.logchan:
			loginfo:=fmt.Sprintf("[%s] [%s] [%s:%s:%d] %s\n", logtmp.timestamp,getloglvstring(logtmp.level),logtmp.filename, logtmp.funcname, logtmp.line,logtmp.msg)
			fmt.Fprint(f.fileobj, loginfo)
			if logtmp.level>=ERROR{
				if f.checksize(f.errfileobj){ //切割日志文件
					newfile,_:=f.splitelog(f.errfileobj)
					f.errfileobj=newfile
				}
				fmt.Fprint(f.errfileobj,loginfo)
			}
		default:
			time.Sleep(time.Millisecond*500)
		}
	
	}
}

//日志分发
func (f *FileLog)logpatch(lv Loglevel,format string,a ...interface{})  {
	if f.enable(lv){
		msg := fmt.Sprintf(format, a...)
		now := time.Now().Format("2006-01-02 15:04:05")
		funcName, fileName, lineNo := getinfo(3)
		//先把日志发送到通道中

		logtmp:=&logmsg{
			level: lv,
			msg: msg,
			funcname: funcName,
			filename: fileName,
			line: lineNo,
			timestamp: now,
		}
		select {
		case f.logchan<-logtmp:
		default:
			//如果日志阻塞,就把日志丢掉，保证不阻塞
		}
	
	}
}



func (f *FileLog) Debug(format string, a ...interface{}) {
	f.logpatch(DEBUG,format ,a...)

}

func (f *FileLog) Trace(format string, a ...interface{}) {
	f.logpatch(TRACE,format ,a...)
}

func (f *FileLog) Info(format string, a ...interface{}) {
	f.logpatch(INFO,format ,a...)
}

func (f *FileLog) Warning(format string, a ...interface{}) {
	f.logpatch(WARNING,format ,a...)
}
func (f *FileLog) Error(format string, a ...interface{}) {
	f.logpatch(ERROR,format ,a...)
}
func (f *FileLog) Fatal(format string, a ...interface{}) {
	f.logpatch(FATAL,format ,a...)
}
