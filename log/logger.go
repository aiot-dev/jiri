// Copyright 2017 The Fuchsia Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"io"
	glog "log"
	"os"
	"sync"
)

// Logger provides for convenient logging in jiri. It supports logger
// level using global flags. To use it "InitializeGlobalLogger" needs to
// be called once, then GetLogger function can be used to get the logger or
// log functions can be called directly
//
// The default logging level is Info. It uses golang logger to log messages internally.
// As an example to use debug logger one needs to run
// log.GetLogger().Debugf(....)
// or
// log.Debugf(....)
// By default Error logger prints to os.Stderr and others print to os.Stdout.
// Capture function can be used to temporarily capture the logs.
type Logger struct {
	lock          *sync.Mutex
	LoggerLevel   LogLevel
	goLogger      *glog.Logger
	goErrorLogger *glog.Logger
}

type LogLevel int

const (
	ErrorLevel LogLevel = iota
	InfoLevel
	DebugLevel
	TraceLevel
	AllLevel
)

func NewLogger(loggerLevel LogLevel) *Logger {
	return &Logger{
		LoggerLevel:   loggerLevel,
		lock:          &sync.Mutex{},
		goLogger:      glog.New(os.Stdout, "", glog.Lmicroseconds),
		goErrorLogger: glog.New(os.Stderr, "", glog.Lmicroseconds),
	}
}

// Capture arranges for the next log to go to supplied io.Writers.
// This will be cleared and not used for any subsequent logs.
// Specifying nil for a writer will result in using the default writer.
// ioutil.Discard should be used to discard output.
func (l Logger) Capture(stdout, stderr io.Writer) Logger {
	if stdout != nil {
		l.goLogger = glog.New(stdout, "", glog.Lmicroseconds)
	}
	if stderr != nil {
		l.goErrorLogger = glog.New(stderr, "", glog.Lmicroseconds)
	}
	return l
}

func (l Logger) log(format string, a ...interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.goLogger.Printf(format, a...)
}

func (l Logger) Infof(format string, a ...interface{}) {
	if l.LoggerLevel >= InfoLevel {
		l.log(format, a...)
	}
}

func (l Logger) Debugf(format string, a ...interface{}) {
	if l.LoggerLevel >= DebugLevel {
		l.log(format, a...)
	}
}

func (l Logger) Tracef(format string, a ...interface{}) {
	if l.LoggerLevel >= TraceLevel {
		l.log(format, a...)
	}
}

func (l Logger) Logf(format string, a ...interface{}) {
	if l.LoggerLevel >= AllLevel {
		l.log(format, a...)
	}
}

func (l Logger) Errorf(format string, a ...interface{}) {
	if l.LoggerLevel >= ErrorLevel {
		l.lock.Lock()
		defer l.lock.Unlock()
		l.goErrorLogger.Printf(format, a...)
	}
}
