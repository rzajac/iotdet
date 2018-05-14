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
    "github.com/pkg/errors"
    jww "github.com/spf13/jwalterweatherman"
    "net"
    "regexp"
)

// iotWiFiRegEx is regular expression matching IoT access point names.
var iotWiFiRegEx *regexp.Regexp = regexp.MustCompile("IOT_([[:xdigit:]]{2}){3}")

// DevAP represents IoT device access point.
// IoT devices create access points when they are waiting to be discovered
// and configured.
type DevAP struct {
    Name   string   // Access point name.
    Bssid  string   // Access point BSSID.
    itf    *wifiItf // The wifi interface we use to connect to this access point.
    stopCh stopChanel
}

// NewDevAP returns new DevAP instance.
func NewDevAP(name, mac string, itf *wifiItf) *DevAP {
    return &DevAP{
        Name:  name,
        Bssid: mac,
        itf:   itf,
    }
}

// IsIotAp checks if access point name matches IoT device.
func (ap *DevAP) IsIotAp() bool {
    return iotWiFiRegEx.MatchString(ap.Name)
}

func (ap *DevAP) Configure(params) error {

}

func (ap *DevAP) connect(pass string) error {
    var err error

    jww.DEBUG.Printf("Connecting to %s with password: %s\n", ap.Name, pass)
    ap.stopCh, err = connectToAp(ap.Name, pass, ap.itf.Name)
    if err != nil {
        return err
    }

    return nil
}

func (ap *DevAP) disconnect() {
    select {
    case ap.stopCh <- struct{}{}:
        jww.DEBUG.Printf("Disconnecting from %s access point.\n", ap.Name)
        <-ap.stopCh
    default:
        return
    }
}

// getIp returns IP address given to WiFi interface.
//
// Note: If interface has more then one IP addresses assigned to it this method will
// return the firs one in the collection.
func (ap *DevAP) getIp() (string, error) {
    var err error
    var ip net.IP
    var addrs []net.Addr

    if addrs, err = ap.itf.Interface.Addrs(); err != nil {
        return "", errors.Wrapf(err, "Can not get IP address for %s.", ap.itf.Name)
    }

    for _, addr := range addrs {
        switch v := addr.(type) {
        case *net.IPNet:
            ip = v.IP
        case *net.IPAddr:
            ip = v.IP
        }

        if ip.To4() != nil {
            return ip.String(), nil
        }
    }

    return "", errors.Errorf("The %s has no IPv4 IP addresses.", ap.itf.Name)
}
