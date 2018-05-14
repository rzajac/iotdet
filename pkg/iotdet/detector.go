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
    "github.com/sirupsen/logrus"
)

// Detector is responsible for detecting new IoT
// devices using specified WiFi interface.
type Detector struct {
    itf *WiFiItf
    log *logrus.Entry
}

// NewDetector returns new Detector instance.
func NewDetector(itfName string, log *logrus.Entry) (*Detector, error) {
    itf, err := GetInterface(itfName, log)
    if err != nil {
        return nil, err
    }

    return &Detector{itf: itf, log: log}, nil
}

// Detect detects IoT access points in range.
func (d *Detector) Detect() ([]*DevAP, error) {
    d.log.Infof("Scanning for new IoT devices using %s interface...", d.itf.Name)
    return d.itf.scanForIotAPs()
}
