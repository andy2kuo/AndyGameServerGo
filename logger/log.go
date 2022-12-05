package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	DEBUG = 0
	INFO  = 1
	WARN  = 2
	ERROR = 3
)

type Logger struct {
	log_filePath string
	log_name     string
	log_env      string

	log_level     int
	rowCount      int
	index         int
	maxRowPerFile int

	logger *log.Logger

	file        *os.File
	createTime  time.Time
	nextLogTime time.Time

	sync.RWMutex
}

func NewLogger(name, env string, level int) *Logger {
	logger := &Logger{log_name: name, log_env: env, log_level: level}
	logger.init()

	return logger
}

func (logger *Logger) init() {
	logger.createTime = time.Now()
	logger.nextLogTime = time.Date(logger.createTime.Year(), logger.createTime.Month(), logger.createTime.Day()-1, 0, 0, 0, 0, logger.createTime.Location())

	logger.rowCount = 0
	logger.index = 0
	logger.maxRowPerFile = 1000

	os.MkdirAll("Log", 0755)
	os.MkdirAll(fmt.Sprintf("Log/%v", logger.log_env), 0755)
	os.MkdirAll(fmt.Sprintf("Log/%v/%v", logger.log_env, logger.log_name), 0755)

	logger.log_filePath = fmt.Sprintf("Log/%v/%v", logger.log_env, logger.log_name)
	logger.refresh()
}

// 
func (logger *Logger) SetMaxRowPerFile(rowCount int) {
	logger.maxRowPerFile = rowCount
}

func (logger *Logger) refresh() {
	if time.Now().Before(logger.nextLogTime) {
		if logger.rowCount >= logger.maxRowPerFile {
			logger.rowCount = 0
			logger.index += 1
		} else {
			return
		}
	} else {
		logger.rowCount = 0
		logger.index = 0
	}

	logger.createTime = time.Now()
	logger.nextLogTime = time.Date(logger.createTime.Year(), logger.createTime.Month(), logger.createTime.Day()+1, 0, 0, 0, 0, logger.createTime.Location())

	newFile, err := os.OpenFile(fmt.Sprintf("%v/%v-%v-%v.log", logger.log_filePath, logger.log_name, logger.createTime.Format("20060102"), logger.index+1), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}

	b, err := io.ReadAll(newFile)
	if err != nil {
		panic(b)
	}

	checkContentList := strings.Split(string(b), "\n")
	nowRowCount := len(checkContentList)
	if len(checkContentList) > 1 {
		logger.rowCount = nowRowCount
	} else {
		logger.rowCount = nowRowCount - 1
	}

	logger.file = newFile
	logger.logger = log.New(logger.file, "", log.LstdFlags)
}

func (logger *Logger) Debug(messages ...interface{}) {
	logger.Lock()
	defer logger.Unlock()

	if DEBUG < logger.log_level {
		return
	}

	logger.refresh()

	logger.logger.SetPrefix("[DEBUG]")
	logger.logger.SetFlags(log.LstdFlags | log.Llongfile)
	logger.logger.Println(messages...)
	logger.rowCount += 1

	fmt.Printf("%v [DEBUG] %v", time.Now().Format("2006/01/02 15:04:05"), fmt.Sprintln(messages...))
}

func (logger *Logger) Info(messages ...interface{}) {
	logger.Lock()
	defer logger.Unlock()

	if INFO < logger.log_level {
		return
	}

	logger.refresh()

	prefix := "[INFO]"
	logger.logger.SetPrefix(prefix)
	logger.logger.SetFlags(log.LstdFlags)
	logger.logger.Println(messages...)
	logger.rowCount += 1

	fmt.Printf("%v [INFO] %v", time.Now().Format("2006/01/02 15:04:05"), fmt.Sprintln(messages...))
}

func (logger *Logger) Warn(messages ...interface{}) {
	logger.Lock()
	defer logger.Unlock()

	if WARN < logger.log_level {
		return
	}

	logger.refresh()

	logger.logger.SetPrefix("[WARN]")
	logger.logger.SetFlags(log.LstdFlags)
	logger.logger.Println(messages...)
	logger.rowCount += 1

	fmt.Printf("%v [WARN] %v", time.Now().Format("2006/01/02 15:04:05"), fmt.Sprintln(messages...))
}

func (logger *Logger) Error(messages ...interface{}) {
	logger.Lock()
	defer logger.Unlock()

	if ERROR < logger.log_level {
		return
	}

	logger.refresh()

	logger.logger.SetPrefix("[ERROR]")
	logger.logger.SetFlags(log.LstdFlags)
	logger.logger.Println(messages...)
	logger.rowCount += 1

	_, file, line, ok := runtime.Caller(1)

	if ok {
		fmt.Printf("%v [ERROR] %v:%v %v", time.Now().Format("2006/01/02 15:04:05"), file, line, fmt.Sprintln(messages...))
	} else {
		fmt.Printf("%v [ERROR] %v", time.Now().Format("2006/01/02 15:04:05"), fmt.Sprintln(messages...))
	}
}

func (logger *Logger) Close() {
	logger.file.Close()
}
