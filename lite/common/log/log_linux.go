package log

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
	"time"
)

func Recover() {
	if runtime.GOOS == "windows" {
		return
	}

	logDirPath := "./log/fatal/"
	os.MkdirAll(logDirPath, 0660)
	files, err := ListDir("./log/fatal/")
	if err == nil {
		for index := range files {
			file, err := os.Stat(logDirPath + files[index])
			if err == nil {
				if file.Size() == 0 {
					os.Remove(logDirPath + files[index])
				}
			}
		}
	}

	logPath := "./log/fatal/" + time.Now().Format("2006-01-02 15:04:05") + ".log"
	recoverFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		fmt.Printf("open panic log file error: %v\n", err)
		return
	}
	syscall.Dup2(int(recoverFile.Fd()), int(os.Stderr.Fd()))

	runtime.SetFinalizer(recoverFile, func(fd *os.File) {
		fd.Close()
	})
}
