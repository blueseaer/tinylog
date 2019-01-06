//logger基本功能：
//1、按天生成文件
//2、当天的文件，超过一定尺寸，进行滚动，生成子文件，文件个数无限制
//3、超过一定天数期限的文件自动删除
//4、仅提供 DEBUG、INFO、ERROR、FATAL四种日志
//5、每种级别日志各自生成一个文件，也可以单独一个文件

package tinylog

import "os"

//输出方式
const (
	WRITE_FILE  = 2 //写文件
	PUT_CONSOLE = 4 //输出到控制台
)

//日志级别
const (
	DEBUG_LEVEL = DEBUG_TYPE
	INFO_LEVEL  = INFO_TYPE
	ERROR_LEVEL = ERROR_TYPE
	FATAL_LEVEL = FATAL_TYPE
)

var logger *TinyLogger

/**********************对外接口**********************/
func Init(name string,
	dir string,
	maxFileSize int64,
	maxDays int,
	level int,
	separate bool,
	outType int32) bool {

	logger = &TinyLogger{}
	logger.level = level
	logger.separate = separate

	// 判断目录是否存在
	s, err := os.Stat(dir)
	if nil != err {
		if nil != os.Mkdir(dir, 0777) {
			return false
		}
	} else {
		if !s.IsDir() {
			return false
		}
	}
	//日志分开，单独一个文件
	if logger.separate {
		for i := DEBUG_TYPE; i <= FATAL_TYPE; i++ {
			log := NewTinylog(name, dir, maxFileSize, maxDays, i, outType)
			log.Init()
			logger.logs[i] = log
		}
	} else { //混合在一个文件
		log := NewTinylog(name, dir, maxFileSize, maxDays, NONE_TYPE, outType)
		log.Init()
		logger.logs[NONE_TYPE] = log
	}
	return true
}

func Debug(format string, a ...interface{}) {
	logger.Write(DEBUG_LEVEL, format, a...)
}

func Info(format string, a ...interface{}) {
	logger.Write(INFO_LEVEL, format, a...)
}

func Error(format string, a ...interface{}) {
	logger.Write(ERROR_LEVEL, format, a...)
}

func Fatal(format string, a ...interface{}) {
	logger.Write(FATAL_LEVEL, format, a...)
}
