// Package clog implements an alternative logger to the one found in the standard
// library with support for more logging levels and a different output format.
// It also has support for splitting log files on daily boundaries.
//
// Author: Clint Caywood
//
// https://github.com/cratonica/clog
package clog

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// Represents how critical the logged
// message is.
type Level uint8

const (
	LevelFatal Level = iota + 1
	LevelError
	LevelWarning
	LevelInfo
	LevelTrace
)

var levelStrings = map[Level]string{
	LevelFatal:   "Fatal",
	LevelError:   "Error",
	LevelWarning: "Warning",
	LevelInfo:    "Info",
	LevelTrace:   "Trace",
}

type output struct {
	writer io.Writer
	level  Level
}

// The Logger
type Clog struct {
	mtx     sync.Mutex
	outputs []output
}

// Instantiate a new Clog
func NewClog() *Clog {
	return &Clog{sync.Mutex{}, make([]output, 0)}
}

// Adds an ouput, specifying the maximum log Level
// you want to be written to this output. For instance,
// if you pass Warning for level, all logs of type
// Warning, Error, and Fatal would be logged to this output.
func (this *Clog) AddOutput(writer io.Writer, level Level) {
	this.outputs = append(this.outputs, output{writer, level})
}

// Convenience function
func (this *Clog) Trace(format string, v ...interface{}) {
	this.Log(LevelTrace, format, v...)
}

// Convenience function
func (this *Clog) Info(format string, v ...interface{}) {
	this.Log(LevelInfo, format, v...)
}

// Convenience function
func (this *Clog) Warning(format string, v ...interface{}) {
	this.Log(LevelWarning, format, v...)
}

// Convenience function
func (this *Clog) Error(format string, v ...interface{}) {
	this.Log(LevelError, format, v...)
}

// Convenience function, will not terminate the program
func (this *Clog) Fatal(format string, v ...interface{}) {
	this.Log(LevelFatal, format, v...)
}

// Logs a message for the given level. Most callers will likely
// prefer to use one of the provided convenience functions.
func (this *Clog) Log(level Level, format string, v ...interface{}) {
	message := fmt.Sprintf(format+"\n", v...)
	strTimestamp := getTimestamp()
	strFinal := fmt.Sprintf("%s [%-7s] %s", strTimestamp, levelStrings[level], message)
	bytes := []byte(strFinal)
	this.mtx.Lock()
	defer this.mtx.Unlock()
	for _, output := range this.outputs {
		if output.level >= level {
			output.writer.Write(bytes)
		}
	}
}

// Gets the timestamp string
func getTimestamp() string {
	now := time.Now()
	return fmt.Sprintf("%v-%02d-%02d %02d:%02d:%02d.%03d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/1000000)
}
