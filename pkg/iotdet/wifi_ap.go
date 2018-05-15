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
    "regexp"
)

// iotWiFiRegEx is regular expression matching IoT access point names.
var iotWiFiRegEx *regexp.Regexp = regexp.MustCompile("IOT_([[:xdigit:]]{2}){3}")

// DevAP represents IoT device access point.
// IoT devices create access points when they are waiting to be discovered
// and configured.
type DevAP struct {
    Name  string   // Access point name.
    Bssid string   // Access point BSSID.
    Itf   *WiFiItf // The wifi interface we use to connect to this access point.
}

// NewDevAP returns new DevAP instance.
func NewDevAP(name, mac string, itf *WiFiItf) *DevAP {
    return &DevAP{
        Name:  name,
        Bssid: mac,
        Itf:   itf,
    }
}

// IsIotAp checks if access point name matches IoT device.
func (ap *DevAP) IsIotAp() bool {
    return iotWiFiRegEx.MatchString(ap.Name)
}

func (ap *DevAP) Connect(pass string) error {
    return ap.Itf.connect(ap.Name, pass)
}

func (ap *DevAP) Disconnect() {
    ap.Itf.disconnect()
}
