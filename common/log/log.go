/*
封装log标准库
*/

package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

//error红色，debug蓝色，info黄色
var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[31m ", log.LstdFlags | log.Lshortfile)
	debugLog = log.New(os.Stdout, "\033[34m[debug]\033[34m ", log.LstdFlags | log.Lshortfile)
	infoLog = log.New(os.Stdout, "\033[33m[info]\033[33m ", log.LstdFlags | log.Lshortfile)
	loggers = []*log.Logger{errorLog, debugLog, infoLog}
	mu sync.Mutex //用于并发控制显示的日志等级
)

//对外暴露不同日志级别的打印方法
var (
	Error = errorLog.Println
	Errorf = errorLog.Printf
	Debug = debugLog.Println
	Debugf = debugLog.Printf
	Info = infoLog.Println
	Infof = infoLog.Printf
)

//设置日志层级
const (
	InfoLevel = iota
	DebugLevel
	ErrorLevel
	Disabled
)

//通过设置等级控制打印，0级别打印所有日志，1级别打印error和debug，以此类推(默认全部打印)
func SetLevel(level int) {

	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if ErrorLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}

	if DebugLevel < level {
		debugLog.SetOutput(ioutil.Discard)
	}

	if InfoLevel < level {
		infoLog.SetOutput(ioutil.Discard)
	}
}