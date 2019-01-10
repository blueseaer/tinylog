## tinylog
Tiny logging package for Golang

## install
go get github.com/idakun/tinylog

## 基本功能
* 按天生成文件
* 当天的文件，超过一定尺寸，进行滚动，生成子文件，文件个数无限制
* 超过一定天数期限的文件自动删除
* 仅提供 DEBUG、INFO、ERROR、FATAL四种日志
* 每种级别日志各自生成一个文件，也可以单独一个文件

## 使用说明
非常简单、易用，只需要一步初始化即可。
```go
package main

import (
	logger "github.com/idakun/tinylog"
)

func main() {
	logger.Init(
		"testlog",                            // 日志文件名前缀
		".",                                  // 日志保存目录
		1024*1024,                            // 文件最大尺寸，单位字节
		3,                                    // 日志保留最大天数
		logger.DEBUG_LEVEL,                   // 日志级别
		true,                                 // 日志是否写在同一个文件
		logger.PUT_CONSOLE|logger.WRITE_FILE) // 日志输出方式

	logger.Debug("This a DEBUG test log.")
	logger.Info("This a INFO test log.")
	logger.Error("This a ERROR test log.")
	logger.Fatal("This a Fatal test log.")
}
```
输出：
```shell
[Debug][2019-01-06 11:33:21.0140][src/main.go:17][main]This a DEBUG test log.
[Info ][2019-01-06 11:33:21.0147][src/main.go:18][main]This a INFO test log.
[Error][2019-01-06 11:33:21.0314][src/main.go:19][main]This a ERROR test log.
[Fatal][2019-01-06 11:33:21.0317][src/main.go:20][main]This a Fatal test log.
```
