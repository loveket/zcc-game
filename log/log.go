package log

import (
	"os"
	"time"
)

var (
	logPath = "G:\\code\\golang\\zcc-game\\log\\logs"
)

type ILog interface {
	ImportantMSG(string)
	WarningMSG(string)
	ErrorMSG(string)
}
type Log struct {
	filepath string
}

func NewLog(filename string) *Log {
	path := logPath + "\\" + filename + ".log"
	return &Log{filepath: path}
}
func (p *Log) ImportantMSG(data string) {
	file, err := os.OpenFile(p.filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer file.Close()
	curTime := time.Now().Format("2006-01-02 15:04:05")
	msgType := "[IMPORT]："
	file.WriteString(curTime + msgType + data + "\n")
}
func (p *Log) WarningMSG(data string) {
	file, err := os.OpenFile(p.filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer file.Close()
	curTime := time.Now().Format("2006-01-02 15:04:05")
	msgType := "[WARNING]："
	file.WriteString(curTime + msgType + data + "\n")
}
func (p *Log) ErrorMSG(data string) {
	file, err := os.OpenFile(p.filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer file.Close()
	curTime := time.Now().Format("2006-01-02 15:04:05")
	msgType := "[ERROR]："
	file.WriteString(curTime + msgType + data + "\n")
}
