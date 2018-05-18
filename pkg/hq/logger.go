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

type Logger interface {
    Debug(args ...interface{})
    Debugf(format string, args ...interface{})

    Info(args ...interface{})
    Infof(format string, args ...interface{})

    Error(args ...interface{})
    Errorf(format string, args ...interface{})
}

type NoopLogger struct{}

func NewNoopLogger() *NoopLogger {
    return &NoopLogger{}
}

func (*NoopLogger) Debug(args ...interface{})                 {}
func (*NoopLogger) Debugf(format string, args ...interface{}) {}
func (*NoopLogger) Info(args ...interface{})                  {}
func (*NoopLogger) Infof(format string, args ...interface{})  {}
func (*NoopLogger) Error(args ...interface{})                 {}
func (*NoopLogger) Errorf(format string, args ...interface{}) {}
