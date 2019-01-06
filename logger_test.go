package tinylog

import (
	"testing"
)

//日志不分离，日志写到文件，同时输出到控制台
func TestLog1(t *testing.T) {
	Init("testlog", ".", 1024*1024, 3, DEBUG_LEVEL, false, PUT_CONSOLE|WRITE_FILE)
	for i := 0; i < 5; i++ {
		Debug("This a DEBUG test log1.")
		Info("This a INFO test log1.")
		Error("This a ERROR test log1.")
		Fatal("This a Fatal test log1.")
	}
}

//日志分离，只输出到控制台
func TestLog2(t *testing.T) {
	Init("testlog", ".", 1024*1024, 3, DEBUG_LEVEL, true, PUT_CONSOLE)
	for i := 0; i < 5; i++ {
		Debug("This a DEBUG test log2.")
		Info("This a INFO test log2.")
		Error("This a ERROR test log2.")
		Fatal("This a Fatal test log2.")
	}
}
