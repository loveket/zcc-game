package log

import (
	"os"
	"time"
	"xiuianserver/utils"
)

var (
	logPath = ""
)

type ILog interface {
	InfoMSG(string)
	WarnMSG(string)
	ErrorMSG(string)
}
type Log struct {
	filepath string
}

func NewLog(filename string) *Log {
	path := utils.GetOsPwd() + "\\log\\diy_logs\\" + filename + ".log"
	return &Log{filepath: path}
}
func (p *Log) InfoMSG(data string) {
	file, err := os.OpenFile(p.filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer file.Close()
	curTime := time.Now().Format("2006-01-02 15:04:05")
	msgType := "[INFO]："
	file.WriteString(curTime + msgType + data + "\n")
}
func (p *Log) WarnMSG(data string) {
	file, err := os.OpenFile(p.filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer file.Close()
	curTime := time.Now().Format("2006-01-02 15:04:05")
	msgType := "[WARN]："
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
