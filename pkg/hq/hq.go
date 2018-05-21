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
    "time"
    "github.com/eclipse/paho.mqtt.golang"
    "fmt"
)

// HQ represents HQ configuration.
type HQ struct {
    // The HQ configuration.
    cfg Configurator
    // The WiFi interface used for new agent detection.
    detItf *wifiItf
    // The MQTT client.
    mqttClient mqtt.Client
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
    log.Infof("scanning for new agents using %s interface...", hq.cfg.GetDetItfName())
    apNames, err := hq.detItf.scan()
    if err != nil {
        return nil, err
    }

    // Filter out non agent access points.
    var agents []*beacon
    for _, apName := range apNames {
        if hq.cfg.GetBeaconNamePat().MatchString(apName) {
            agents = append(agents, hq.agentAP(apName))
        }
    }

    return agents, nil
}

// Configure configure given agent access point.
func (hq *HQ) Configure(apName string) error {
    ap := hq.agentAP(apName)

    cmd := hq.getConfigCmd().MarshalCmd()
    resp, err := hq.detItf.sendCmd(ap, cmd)
    if err != nil {
        return err
    }

    log.Infof("agent response: %s", string(resp))

    return nil
}

// SetMQTTConfig sets MQTT broker configuration.
func (hq *HQ) SetMQTTConfig(cfg MQTTConfig) error {
    opts := mqtt.NewClientOptions()
    opts.AddBroker(fmt.Sprintf("tcp://%s:%d", cfg.ip, cfg.port))
    opts.SetClientID(cfg.clientID)
    opts.SetKeepAlive(2 * time.Second)
    opts.SetPingTimeout(1 * time.Second)
    opts.SetUsername(cfg.user)
    opts.SetPassword(cfg.pass)

    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        return token.Error()
    }

    hq.mqttClient = client
    return nil
}

// PublishMQTT publishes payload to MQTT broker.
func (hq *HQ) PublishMQTT(topic string, qos byte, retained bool, payload interface{}) error {
    hq.mqttClient.Publish(topic, qos, retained, payload).Wait()
    return nil
}

// getConfigCmd returns configuration command.
func (hq *HQ) getConfigCmd() *cmdConfig {
    cmd := newConfigCmd()
    cmd.ApName = hq.cfg.GetAPName()
    cmd.ApPass = hq.cfg.GetAPPass()

    if hq.IsMQTTSet() {
        cmd.MQTTIP = hq.cfg.GetMQTTIP()
        cmd.MQTTPort = hq.cfg.GetMQTTPort()
        cmd.MQTTUser = hq.cfg.GetMQTTUser()
        cmd.MQTTPass = hq.cfg.GetMQTTPass()
    }

    return cmd
}

// IsMQTTSet returns true if MQTT configuration has been set.
func (hq *HQ) IsMQTTSet() bool {
    return hq.cfg.GetMQTTIP() != "" && hq.cfg.GetMQTTPort() != 0
}

// beacon is helper method returning configured
// agent access point with given name.
func (hq *HQ) agentAP(apName string) *beacon {
    return newBeacon(apName,
        hq.cfg.GetBeaconPass(),
        hq.cfg.GetBeaconIP(),
        hq.cfg.GetDetItfIP(),
        hq.cfg.GetBeaconPort())
}
