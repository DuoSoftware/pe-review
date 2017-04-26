package Common

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

var fileHandlers map[string]*os.File
var fileHandlerLock = sync.RWMutex{}

func GetFileHandler(handler string) (conn *os.File) {
	fileHandlerLock.RLock()
	defer fileHandlerLock.RUnlock()
	conn = fileHandlers[handler]
	return
}

func SetFileHandler(handler string, conn *os.File) {
	fileHandlerLock.Lock()
	defer fileHandlerLock.Unlock()
	fileHandlers[handler] = conn
}

func PublishLog(fileName string, Body string) {
	if fileHandlers == nil {
		fileHandlers = make(map[string]*os.File)
	}
	var slash = ""
	if runtime.GOOS == "windows" {
		slash = "\\"
	} else {
		slash = "/"
	}
	var logfolderName = "Logs"

	// check if the publishing activity folder exists, if not create one.

	logrootpath := logfolderName
	_, activityRooterr := os.Stat(logrootpath)
	if activityRooterr != nil {
		// create folder in the given path and permissions
		os.Mkdir(logrootpath, 0777)
	}

	var ff *os.File
	var err error
	name := logfolderName + slash + fileName
	if GetFileHandler(name) != nil {
		ff = GetFileHandler(name)
	} else {
		ff, err = os.OpenFile(name, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			ff, err = os.Create(name)
			if err == nil {
				ff, err = os.OpenFile(name, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
				if err == nil {
					SetFileHandler(name, ff)
				}
			}
		}

	}

	_, err = ff.Write([]byte(string(time.Now().Local().Format("2006-01-02 @ 15:04:05")) + "  " + Body + "  " + "\r\n"))
	if err != nil {
		fmt.Println(err.Error())
	}

}

func removeFile(filename string) {
	err := os.Remove(filename)
	if err != nil {
		fmt.Println(err)
	}
}
