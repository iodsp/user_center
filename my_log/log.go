package my_log

import (
	"fmt"
	"github.com/iodsp/user_center/context"
	"github.com/op/go-logging"
	"os"
	"time"
)

type Log struct {
	Logger *logging.Logger
}

func NewLog(conf *context.Config) *Log {
	logFile := LogFilePath()
	backend1 := logging.NewLogBackend(logFile, "", 0)
	var format1 = logging.MustStringFormatter(
		`%{level:.4s} %{time:2006-01-02T15:04:05.999}  %{id:03x} %{message}`,
	)
	backend1Leveled := logging.NewBackendFormatter(backend1, format1)

	if conf.Debug() {
		backend2 := logging.NewLogBackend(os.Stderr, "", 0)
		var format2 = logging.MustStringFormatter(
			`%{color}%{time:2006-01-02T15:04:05.999} %{level:.4s} %{id:03x}%{color:reset} %{message}`,
		)
		backend2Formatter := logging.NewBackendFormatter(backend2, format2)
		logging.SetBackend(backend2Formatter, backend1Leveled)
	} else {
		logging.SetBackend(backend1Leveled)
	}

	instance := &Log{
		Logger: logging.MustGetLogger("example"),
	}
	return instance
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func LogFilePath() *os.File {
	currentDate := time.Now().Format("2006-01-02")
	path := "logs/"
	isExist, _ := pathExists(path)
	if isExist == false {
		dirErr := os.MkdirAll(path, 0777)
		if nil != dirErr {
			fmt.Println(dirErr)
		}
	}
	fileName := path + "/" + currentDate + ".log"
	logFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("fail to open log file " + fileName)
	}
	return logFile
}

//todo 拿到每个接口的入参出参 记录到日志
