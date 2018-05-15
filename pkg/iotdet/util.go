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

package iotdet

import (
    "time"
    "fmt"
    "os"
)

type stopChanel chan struct{}
type connChanel chan struct{}

type CtrlChanel chan string

const CTRL_STOP = "STOP"

// runUntil runs given function until it returns true or atMost times tries.
func runUntil(what func() bool, interval time.Duration, atMost int) chan bool {
    var try int = 0
    stop := make(chan bool)

    go func() {
        for {
            try += 1
            if what() {
                stop <- true
                break
            }
            select {
            case <-time.After(interval):
            }
            if try == atMost {
                stop <- false
                break
            }
        }
    }()

    return stop
}

func dumpBytes(buf []byte) {
    for i, x := range buf {
        fmt.Printf("%02x ", x)
        if (i+1)%8 == 0 {
            fmt.Print(" ")
        }
        if (i+1)%16 == 0 {
            fmt.Println("")
        }
    }
    fmt.Print("\n\n")
}

// Checks if directory exists.
func dirExists(dirPath string) bool {
    if _, err := os.Stat(dirPath); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}
