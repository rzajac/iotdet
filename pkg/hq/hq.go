// HomeHQ (c) 2018 Rafal Zajac <rzajac@gmail.com> All rights reserved.
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

// HQ represents HQ configuration.
type HQ struct {
    // The HQ configuration.
    cfg Configurator
    // The WiFi interface used for new agent detection.
    detItf *wifiItf
}

// NewHQ returns new instance of HQ with some default values for fields.
func NewHQ(cfg Configurator) (*HQ, error) {
    itf, err := getInterface(cfg.GetDetItfName())
    if err != nil {
        return nil, err
    }
    h := &HQ{
        cfg:    cfg,
        detItf: itf,
    }
    return h, nil
}

// Detect detects IoT access points in range.
func (hq *HQ) DetectAgents() ([]*beacon, error) {
    INFO.Printf("scanning for new agents using %s interface\n", hq.cfg.GetDetItfName())
    apNames, err := hq.detItf.scan()
    if err != nil {
        return nil, err
    }

    // Filter out non agent access points.
    var agents []*beacon
    for _, apName := range apNames {
        if hq.cfg.GetBeaconNamePat().MatchString(apName) {
            agents = append(agents, hq.beacon(apName))
        }
    }

    return agents, nil
}

// Configure configure given agent access point.
func (hq *HQ) Configure(apName string) error {
    ap := hq.beacon(apName)

    cmd := hq.getConfigCmd().MarshalCmd()
    resp, err := hq.detItf.sendCmd(ap, cmd)
    if err != nil {
        return err
    }

    INFO.Println("agent response: ", string(resp))

    return nil
}

// getConfigCmd returns configuration command.
func (hq *HQ) getConfigCmd() *cmdConfig {
    cmd := newConfigCmd()
    cmd.ApName = hq.cfg.GetAPName()
    cmd.ApPass = hq.cfg.GetAPPass()
    cmd.MQTTIP = hq.cfg.GetMQTTIP()
    cmd.MQTTPort = hq.cfg.GetMQTTPort()
    cmd.MQTTUser = hq.cfg.GetMQTTUser()
    cmd.MQTTPass = hq.cfg.GetMQTTPass()

    return cmd
}

// beacon is helper method returning configured beacon instance.
func (hq *HQ) beacon(name string) *beacon {
    return newBeacon(name,
        hq.cfg.GetBeaconPass(),
        hq.cfg.GetBeaconIP(),
        hq.cfg.GetItfIP(),
        hq.cfg.GetBeaconPort())
}
