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
    Debug(args ...interface{})
    Debugf(format string, args ...interface{})
    Info(args ...interface{})
    Infof(format string, args ...interface{})
    Error(args ...interface{})
    Errorf(format string, args ...interface{})
}

type Log struct {
    debug bool
}

func NewLog() *Log {
    return &Log{}
}

func (l *Log) DebugOn() *Log {
    l.debug = true
    return l
}

func (l *Log) Debug(args ...interface{}) {
    if !l.debug {
        return
    }
    fmt.Fprintf(os.Stdout, "DEBUG: "+fmt.Sprint(args...))
}

func (l *Log) Debugf(format string, args ...interface{}) {
    if !l.debug {
        return
    }
    fmt.Fprintf(os.Stdout, "DEBUG: "+format, args...)
}

func (*Log) Info(args ...interface{}) {
    fmt.Fprintf(os.Stdout, "INFO: "+fmt.Sprint(args...))
}

func (*Log) Infof(format string, args ...interface{}) {
    fmt.Fprintf(os.Stdout, "INFO: "+format, args...)
}

func (*Log) Error(args ...interface{}) {
    fmt.Fprintf(os.Stderr, "ERROR: "+fmt.Sprint(args...))
}

func (*Log) Errorf(format string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, "ERROR: "+format, args...)
}
