// IoTDet (c) 2017 Rafal Zajac <rzajac@gmail.com> All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hq

import (
    "fmt"
    "os"
)

// Logger is an interface for logging.
type Logger interface {
    // Debug logs a message at level Debug on the standard output.
    Debug(args ...interface{})
    // Debugf logs a message at level Debug on the standard output.
    Debugf(format string, args ...interface{})
    // Info logs a message at level Info on the standard output.
    Info(args ...interface{})
    // Infof logs a message at level Info on the standard output.
    Infof(format string, args ...interface{})
    // Error logs a message at level Error on the standard error.
    Error(args ...interface{})
    // Errorf logs a message at level Error on the standard error.
    Errorf(format string, args ...interface{})
}

// Log represents a logger.
type Log struct {
    debug bool
}

// NewLog creates new instance of logger.
func NewLog() *Log {
    return &Log{}
}

// DebugOn turns on debug level logging.
func (l *Log) DebugOn() *Log {
    l.debug = true
    return l
}

func (l *Log) Debug(args ...interface{}) {
    if !l.debug {
        return
    }
    fmt.Fprintf(os.Stdout, "DEBUG: "+fmt.Sprint(args...)+"\n")
}

func (l *Log) Debugf(format string, args ...interface{}) {
    if !l.debug {
        return
    }
    fmt.Fprintf(os.Stdout, "DEBUG: "+format+"\n", args...)
}

func (*Log) Info(args ...interface{}) {
    fmt.Fprintf(os.Stdout, "INFO: "+fmt.Sprint(args...)+"\n")
}

func (*Log) Infof(format string, args ...interface{}) {
    fmt.Fprintf(os.Stdout, "INFO: "+format+"\n", args...)
}

func (*Log) Error(args ...interface{}) {
    fmt.Fprintf(os.Stderr, "ERROR: "+fmt.Sprint(args...)+"\n")
}

func (*Log) Errorf(format string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
}
