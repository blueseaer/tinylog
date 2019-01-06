package tinylog

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"
	"time"
)

const (
	NONE_TYPE = iota //表示无类型，所以日志皆可以写入
	DEBUG_TYPE
	INFO_TYPE
	ERROR_TYPE
	FATAL_TYPE
)
const (
	MAX_LOG_TYPE = 5
)

var LevelNames = []string{"Log", "Debug", "Info ", "Error", "Fatal"}

type TinyLogger struct {
	logs     [MAX_LOG_TYPE]*TinyLog
	level    int
	separate bool
}

func (logger *TinyLogger) Write(level int, format string, a ...interface{}) {
	if level < logger.level {
		return
	}
	if logger.separate {
		logger.logs[level].Write(level, format, a...)
	} else {
		logger.logs[NONE_TYPE].Write(level, format, a...)
	}
}

type TinyLog struct {
	BaseName        string   //日志文件基本名称
	Dir             string   //所在的目录
	MaxFileSize     int64    //文件最大尺寸
	MaxDays         int      //文件最大保留的天数
	OutType         int32    //日志输出方式，目前只支持文件和控制台
	logType         int      //日志类型 只有5种
	curFileHandle   *os.File //当前文件句柄
	curFilePath     string   //当前文件路径
	curSubFileCount int      //当前子文件个数
	mutex           sync.Mutex
}

func NewTinylog(name string, dir string, maxFileSize int64, maxDays int, logType int, outType int32) *TinyLog {
	return &TinyLog{
		BaseName:    name,
		Dir:         dir,
		MaxFileSize: maxFileSize,
		MaxDays:     maxDays,
		OutType:     outType,
		logType:     logType,
	}
}

func (log *TinyLog) Init() {
	if log.OutType&WRITE_FILE != WRITE_FILE {
		// 不用创建文件
		return
	}
	filePath := log.GetFilePathByTime(time.Now())
	h, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if nil != err {
		panic(fmt.Sprintf("OpenFile(%s) failed, error:%s", filePath, err.Error()))
	}
	if nil != log.curFileHandle {
		log.curFileHandle.Close()
		log.curFileHandle = nil
	}
	log.curFileHandle = h
	log.curFilePath = filePath
}

func (log *TinyLog) Write(level int, format string, a ...interface{}) {
	log.mutex.Lock()
	defer log.mutex.Unlock()
	//滚动日志
	log.Rotate()
	pc, file, line, _ := runtime.Caller(3)
	now := time.Now().Format("2006-01-02 15:04:05.0000")
	//[级别][时间戳][文件名:行号][函数名]
	prefix := fmt.Sprintf("[%s][%s][%s:%d][%s]",
		LevelNames[level],
		now,
		log.GetShortName('/', file, 2),
		line,
		log.GetShortName('.', runtime.FuncForPC(pc).Name(), 1))
	//用户日志
	userLog := fmt.Sprintf(format+"\n", a...)

	if log.OutType&WRITE_FILE == WRITE_FILE {
		_, err := log.curFileHandle.Write([]byte(prefix + userLog))
		if nil != err {
			panic("Write failed, error:" + err.Error())
		}
		log.curFileHandle.Sync()
	}
	if log.OutType&PUT_CONSOLE == PUT_CONSOLE {
		fmt.Print(prefix + userLog)
	}

}

func (log *TinyLog) Rotate() {
	if log.GetFilePathByTime(time.Now()) == log.curFilePath {
		log.SizeRoll()
	} else {
		log.TimeRoll()
	}
}

func (log *TinyLog) SizeRoll() {
	info, err := os.Stat(log.curFilePath)
	//有可能文件被删掉
	if err != nil {
		log.Init()
		return
	}
	if info.Size() < log.MaxFileSize {
		return
	}
	if 0 == log.curSubFileCount {
		log.FindSubFileCount()
	}
	//文件分割
	for i := log.curSubFileCount; i >= 0; i-- {
		newFilePath := fmt.Sprintf("%s.%d", log.curFilePath, i+1)
		if 0 == i {
			//先关闭当前文件，后Rename，再初始化
			log.curFileHandle.Close()
			log.curFileHandle = nil
			os.Rename(log.curFilePath, newFilePath)
			log.Init()
		} else {
			oldFilePath := fmt.Sprintf("%s.%d", log.curFilePath, i)
			os.Rename(oldFilePath, newFilePath)
		}
	}
	log.curSubFileCount++
}

func (log *TinyLog) TimeRoll() {
	log.Init()

	now := time.Now().Unix()
	beforeUnix := now - int64(log.MaxDays*24*3600)
	beforeTime := time.Unix(beforeUnix, 0)
	beforeFilePath := log.GetFilePathByTime(beforeTime)

	//执行shell，删除MaxDays天之前的日志
	go func() {
		removeCmd := "rm -rf " + beforeFilePath + "*"
		cmd := exec.Command("sh", "-c", removeCmd)
		cmd.Output()
	}()
}

func (log *TinyLog) GetFilePathByTime(t time.Time) string {
	//文件名格式:Server20181024.debug
	fileName := log.BaseName + t.Format("20060102") + "." + LevelNames[log.logType]
	return log.Dir + "/" + fileName
}

//home/KentZhang/golang/src/network/agent.go 只需要network/agent.go即可
func (log *TinyLog) GetShortName(char byte, name string, pos int) string {
	path := []byte(name)
	count := 0
	i := len(path) - 1
	for ; i >= 0; i-- {
		if char == path[i] {
			count++
		}
		if count >= pos {
			break
		}
	}
	return string(path[i+1:])
}

//如果进程重启，由于curSubFileCount值丢失，再次SizeRoll时，会造成日志丢失，这里要通过遍历文件找回curSubFileCount的值
func (log *TinyLog) FindSubFileCount() {
	proc := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}
		exp := fmt.Sprintf(`^%s.[\d]`,
			log.GetShortName('/', log.curFilePath, 1))
		filepathReg := regexp.MustCompile(exp)
		matchs := filepathReg.FindAllString(path, 1)
		if len(matchs) > 0 {
			log.curSubFileCount++
		}
		return nil
	}
	filepath.Walk(log.Dir, proc)
}
