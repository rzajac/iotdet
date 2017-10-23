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
    jww "github.com/spf13/jwalterweatherman"
)

// Detect detects IoT access points in range.
func Detect(cfg *IotCfg) ([]*AccessPoint, error) {
    var err error
    var itf *wifiItf

    if itf, err = findWifiInterface(cfg.ItfName); err != nil {
        return nil, err
    }

    jww.INFO.Println("Scanning for new IoT devices...")
    return itf.scanForIotAPs()
}

// Configure configures IoT devices.
func Configure(cfg *IotCfg, aps []*AccessPoint) error {
    var err error

    var iotDev *iotDev
    for _, ap := range aps {
        if err = ap.connect(cfg.DevApPass); err != nil {
            jww.ERROR.Println(err)
            ap.disconnect()
            continue
        }

        if err = setIp(cfg.ItfName, cfg.MyIp); err != nil {
            ap.disconnect()
            return err
        }

        if err = pingIot(cfg.ItfName, cfg.IotIp); err != nil {
            ap.disconnect()
            return err
        }

        iotDev = newIotDev(cfg.IotIp, cfg.Cipher)
        if _, err = iotDev.sendCmd(cfg.TcpPort, newApCmd(cfg.IotApName, cfg.IotApPass)); err != nil {
            jww.ERROR.Println(err)
        }

        ap.disconnect()
    }

    return nil
}
