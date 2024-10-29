package log

import (
	"fmt"
	"github.com/utahta/go-cronowriter"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	errorRed   = 31
	warnYellow = 33
	infoBlue   = 34
	debugBlack = 35
)

const (
	LevelError = errorRed
	LevelWarn  = warnYellow
	LevelInfo  = infoBlue
	LevelDebug = debugBlack
)

type Log struct {
	w           *cronowriter.CronoWriter
	preFit      string
	write2File  bool
	write2Front bool
	debugLog    bool
}

var defaultLog *Log

func setColor(msg string, text int) string {
	return fmt.Sprintf("%c[%dm%s%c[0m", 0x1B, text, msg, 0x1B)
}

func (this *Log) Info(format string, v ...interface{}) {
	this.PrintLog(infoBlue, fmt.Sprintf(format, v...))
}

func (this *Log) Error(format string, v ...interface{}) {
	this.PrintLog(errorRed, fmt.Sprintf(format, v...))
}

func (this *Log) Warn(format string, v ...interface{}) {
	this.PrintLog(warnYellow, fmt.Sprintf(format, v...))
}

func (this *Log) Debug(format string, v ...interface{}) {
	this.PrintLog(debugBlack, fmt.Sprintf(format, v...))
}

func (this *Log) PrintLog(lev int, data string) {
	_, file, line, ok := runtime.Caller(2)

	var log = "["
	if this.preFit != "" {
		log = this.preFit + " "
	}
	location, _ := time.LoadLocation("Asia/Shanghai")
	log += time.Now().In(location).Format("2006-01-02 15:04:05") + " "
	panding := 19
	if ok {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = strings.ReplaceAll(file[i+1:], ".go", "")
				break
			}
		}
		if len(short) >= 15 {
			short = short[:15]
			panding -= 15
		} else {
			panding -= len(short)
		}
		log += short
		log += ":"
		_line := strconv.Itoa(line)
		log += _line
		panding -= len(_line)
	}

	log += strings.Repeat(" ", panding)
	log += "]\t"
	switch lev {
	case errorRed:
		log += "Error"
		break
	case warnYellow:
		log += "Warn"
		break
	case infoBlue:
		log += "Info"
		break
	case debugBlack:
		log += "Debug"
		break
	}
	log += ": "
	log += data
	log += "\n"

	log = setColor(log, lev)
	if this.write2File {
		this.w.Write([]byte(log))
	}
	if this.write2Front {
		if len(log) > 250 {
			log = log[0:250] + "\n"
		}
		fmt.Print(log)
	}
}

func init() {
	if defaultLog == nil {
		InitDefaultLog("", true, false)
	}
}

func Debug(format string, v ...interface{}) {
	if defaultLog.debugLog {
		defaultLog.PrintLog(debugBlack, fmt.Sprintf(format, v...))
	}
}

func Info(format string, v ...interface{}) {
	defaultLog.PrintLog(infoBlue, fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	defaultLog.PrintLog(errorRed, fmt.Sprintf(format, v...))
}

func Warn(format string, v ...interface{}) {
	defaultLog.PrintLog(warnYellow, fmt.Sprintf(format, v...))
}

func Logs(leave int, format string, v ...interface{}) {
	defaultLog.PrintLog(warnYellow, fmt.Sprintf(format, v...))
}

func ListDir(dirPth string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		files = append(files, fi.Name())
	}
	return files, nil
}

func getTimeArr(start, end string) int64 {
	timeLayout := "20060102"
	loc, _ := time.LoadLocation("Local")
	startUnix, _ := time.ParseInLocation(timeLayout, start, loc)
	endUnix, _ := time.ParseInLocation(timeLayout, end, loc)
	startTime := startUnix.Unix()
	endTime := endUnix.Unix()
	date := (endTime - startTime) / 86400
	return date
}

func getLogTime(fileName string) string {
	if len(fileName) < 5 {
		return ""
	}
	p := strings.Index(fileName[4:], "_")
	if p == -1 {
		return ""
	}
	p += 5
	if len(fileName) < p+8 {
		return ""
	}
	return fileName[p : p+8]
}

func CleanLog() {
	path, _ := os.Getwd()
	today := time.Now().Format("20060102")
	logs, _ := ListDir(path + "/log")
	for _, logFile := range logs {
		logTime := getLogTime(logFile)
		if logTime == "" {
			continue
		}

		arr := getTimeArr(logTime, today)
		if arr >= 3 {
			os.Remove(path + "/log/" + logFile)
		}
	}
}

func cleanLog(t1 interface{}) {
	for {
		select {
		case <-t1.(*time.Ticker).C:
			CleanLog()
		}
	}
}

func TestPanic() {
	var t []byte
	t[10] = 0
}

var recoverFile *os.File

func InitDefaultLog(logName string, write2Front bool, write2File bool) {
	Recover()
	CleanLog()
	path, _ := os.Getwd()
	defaultLog = &Log{}
	defaultLog.preFit = ""
	defaultLog.write2Front = write2Front
	defaultLog.write2File = write2File
	defaultLog.w = cronowriter.MustNew(path + "/log/log_" + logName + "_%Y%m%d.txt")

	timeClean := time.NewTicker(time.Hour * 24)
	go cleanLog(timeClean)
}

func NewFileLog(name string) *Log {
	log := &Log{}
	log.preFit = ""
	log.write2Front = false
	log.write2File = true
	log.w = cronowriter.MustNew("./" + name + "_%Y%m%d.txt")
	return log
}

func DisAbleDebugLog() {
	defaultLog.debugLog = false
}

func EnAbleDebugLog() {
	defaultLog.debugLog = true
}
