// Package clog implements an alternative logger to the one found in the standard
// library with support for more logging levels and a different output format.
// This is not exhaustive or feature-rich.
//
// Author: Clint Caywood (www.clintcaywood.com)
package clog

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type Level uint8

const (
	LevelFatal Level = iota + 1
	LevelError
	LevelWarning
	LevelInfo
	LevelTrace
)

var LevelStrings = map[Level]string{
	LevelFatal:   "Fatal",
	LevelError:   "Error",
	LevelWarning: "Warning",
	LevelInfo:    "Info",
	LevelTrace:   "Trace",
}

type Output struct {
	writer io.Writer
	level  Level
}

type Clog struct {
	mtx     sync.Mutex
	outputs []Output
}

func NewClog() *Clog {
	return &Clog{sync.Mutex{}, make([]Output, 0)}
}

func (this *Clog) AddOutput(writer io.Writer, level Level) {
	this.outputs = append(this.outputs, Output{writer, level})
}

func (this *Clog) Trace(format string, v ...interface{}) {
	this.Log(LevelTrace, format, v...)
}

func (this *Clog) Info(format string, v ...interface{}) {
	this.Log(LevelInfo, format, v...)
}

func (this *Clog) Warning(format string, v ...interface{}) {
	this.Log(LevelWarning, format, v...)
}

func (this *Clog) Error(format string, v ...interface{}) {
	this.Log(LevelError, format, v...)
}

// Will not terminate the program
func (this *Clog) Fatal(format string, v ...interface{}) {
	this.Log(LevelFatal, format, v...)
}

// Logs a message
func (this *Clog) Log(level Level, format string, v ...interface{}) {
	message := fmt.Sprintf(format+"\n", v...)
	strTimestamp := getTimestamp()
	strFinal := fmt.Sprintf("%s [%-7s] %s", strTimestamp, LevelStrings[level], message)
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
