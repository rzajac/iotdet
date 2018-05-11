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
    "time"
)

func findWifiInterface(itfName string) (*wifiItf, error) {
    var err error
    var itfs []net.Interface

    if itfs, err = getWifiInterfaces(); err != nil {
        return nil, err
    }

    // Check if given interface name exists on the system.
    for _, itf := range itfs {
        if itf.Name == itfName {
            return &wifiItf{itf}, nil
        }
    }

    return nil, errors.Errorf("wifi interface %s not found", itfName)
}

type wifiItf struct {
    net.Interface
}

// makeSureIsUp makes sure interface is up.
func (wi *wifiItf) makeSureIsUp() error {
    if isUp(wi.Interface) {
        return nil
    }

    jww.DEBUG.Printf("Waiting for %s wifi interface to became available.", wi.Name)
    if err := ifUp(wi.Interface); err != nil {
        return err
    }

    stopCh := runUntil(func() bool {
        return isUp(wi.Interface)
    }, 1*time.Second, 5)

    success := <-stopCh
    if !success {
        return errors.Errorf("interface %s could not be brought up", wi.Name)
    }

    jww.DEBUG.Printf("wifiItf %s is up.\n", wi.Name)

    return nil
}

// scanForAPs returns a list of WiFi access points in range.
func (wi *wifiItf) scanForAPs() ([]*AccessPoint, error) {
    var aps []*AccessPoint
    if err := wi.makeSureIsUp(); err != nil {
        return aps, err
    }

    return scanForAPs(wi)
}

// scanForIotAPs returns a list of IoT access points in range.
func (wi *wifiItf) scanForIotAPs() ([]*AccessPoint, error) {
    var aps []*AccessPoint
    if err := wi.makeSureIsUp(); err != nil {
        return aps, err
    }

    aps, err := scanForAPs(wi)
    if err != nil {
        return nil, err
    }

    var iotAps []*AccessPoint
    for _, ap := range aps {
        jww.INFO.Printf("Found AP %s (%s).\n", ap.Name, ap.Bssid)
        if ap.IsIotAp() {
            iotAps = append(iotAps, ap)
        }
    }

    return iotAps, nil
}
