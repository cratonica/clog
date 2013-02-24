package clog

import (
    "os"
    "time"
    "sync"
    "fmt"
    "log"
)

// A File that automatically gets split on daily boundaries. 
// Implements io.Writer.
type DailyFile struct {
    currentFile *os.File
    filenameFormat string
    lastFileChange time.Time
    mtx sync.Mutex
}

// Creates a new DailyFile using the specified file name format.
// Format needs to have one %s that is replaced by the current
// day's date in the format of YYYY-MM-DD. For instance,
// passing "/opt/logs/myprogram_%s.log" would create
// log files like "/opt/logs/myprogram_2013-02-24.log", 
// "/opt/logs/myprogram_2013-02-25.log". Days where nothing
// is written will be skipped.
func NewDailyFile(filenameFormat string) *DailyFile {
    return &DailyFile{nil, filenameFormat, time.Unix(0, 0), sync.Mutex{}}
}

func (this *DailyFile) Write(p []byte) (n int, err error) {
    this.mtx.Lock();
    defer this.mtx.Unlock()
    now := time.Now()
    if this.currentFile == nil || now.Day() != this.lastFileChange.Day() || now.Sub(this.lastFileChange).Hours() > 24 {
        if err := this.rollToNextFile(now); err != nil {
            return 0, err
        }
    }
    return this.currentFile.Write(p)
}

func (this *DailyFile) rollToNextFile(now time.Time) error {
    if (this.currentFile != nil) {
        this.currentFile.Close()
    }
    strDate := fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())
    newFileName := fmt.Sprintf(this.filenameFormat, strDate)
    var err error
    this.currentFile, err = os.OpenFile(newFileName, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0666)
    if err != nil {
        log.Printf("Clog: Unable to open %v for writing\n", newFileName)
        return err
    }
    return err
}

