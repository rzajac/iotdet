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
    "net"
    "time"
    "sync"
    "github.com/sirupsen/logrus"
)

// interfaces is a collection of system WiFi interfaces.
var interfaces map[string]*wifiItf

func init() {
    interfaces = make(map[string]*wifiItf, 2)
}

// GetInterface returns WiFi interface or error if it does not exist.
func GetInterface(name string, log *logrus.Entry) (*wifiItf, error) {
    if itf, ok := interfaces[name]; ok {
        return itf, nil
    }

    var err error
    var itfs []net.Interface
    if itfs, err = getWifiInterfaces(); err != nil {
        return nil, err
    }

    // Check if given interface name exists on the system.
    for _, itf := range itfs {
        if itf.Name == name {
            interfaces[name] = &wifiItf{
                Interface: itf,
                Mutex:     &sync.Mutex{},
                log:       log,
            }
            return interfaces[name], nil
        }
    }

    return nil, errors.Errorf("wifi interface %s not found", name)
}

// wifiItf represents WiFi interface.
type wifiItf struct {
    net.Interface
    *sync.Mutex
    log    *logrus.Entry
    discCh stopChanel
}

// SetIP sets IP on the interface.
func (w *wifiItf) SetIP(ip string) error {
    return setIp(w.Name, ip)
}

// Ping pings IP address and returns error if IP cannot be pinged.
func (w *wifiItf) Ping(ip string) error {
    return ping(w.Name, ip)
}

// Connect connects to access point.
func (w *wifiItf) Connect(apName, apPass string) error {
    var err error
    w.discCh, err = connectToAp(apName, apPass, w.Name)
    if err != nil {
        return err
    }

    return nil
}

// Disconnect disconnects from access point.
func (w *wifiItf) Disconnect() {
    select {
    case w.discCh <- struct{}{}:
        w.log.Debugf("disconnecting %s", w.Name)
        <-w.discCh
    default:
        return
    }
}

// makeSureIsUp makes sure interface is up.
func (w *wifiItf) makeSureIsUp() error {
    if isUp(w.Interface) {
        return nil
    }

    w.log.Debugf("waiting for %s wifi interface to became available", w.Name)
    if err := ifUp(w.Interface); err != nil {
        return err
    }

    stopCh := runUntil(func() bool {
        return isUp(w.Interface)
    }, 1*time.Second, 5)

    success := <-stopCh
    if !success {
        return errors.Errorf("interface %s could not be brought up", w.Name)
    }

    w.log.Debugf("wifiItf %s is up.\n", w.Name)

    return nil
}

// scanForAPs returns a list of WiFi access points in range.
func (w *wifiItf) scanForAPs() ([]*DevAP, error) {
    var aps []*DevAP
    if err := w.makeSureIsUp(); err != nil {
        return aps, err
    }

    return scanForAPs(w)
}

// scanForIotAPs returns a list of IoT access points in range.
func (w *wifiItf) scanForIotAPs() ([]*DevAP, error) {
    var aps []*DevAP
    if err := w.makeSureIsUp(); err != nil {
        return aps, err
    }

    aps, err := scanForAPs(w)
    if err != nil {
        return nil, err
    }

    var iotAps []*DevAP
    for _, ap := range aps {
        w.log.Infof("found AP %s (%s)", ap.Name, ap.Bssid)
        if ap.IsIotAp() {
            iotAps = append(iotAps, ap)
        }
    }

    return iotAps, nil
}
