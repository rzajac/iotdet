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
    "github.com/pkg/errors"
    "regexp"
    "github.com/eclipse/paho.mqtt.golang"
    "fmt"
)

// HQ represents HQ configuration.
type HQ struct {
    // The interface name to use for new agents detection.
    detItfName string
    // The WiFi interface used for new agent detection.
    detItf *wifiItf
    // The regexp used to recognize agent's access point names.
    detApNamePat *regexp.Regexp
    // The password for agent's access point.
    detApPass string
    // The IP agent assigns to itself when creating access point.
    detAgentIP string
    // The IP to use after connecting to agent's access point.
    detUseIP string
    // The TCP port agents listen on for configuration commands.
    detCmdPort int
    // The interval for scanning for new agents.
    detInterval time.Duration
    // The cipher to use for communication with agent devices.
    cipher Cipher
    // The access point name agents use to communicate.
    mainApName string
    // The access point password.
    mainApPass string
    // Logger configuration.
    log Logger
    // The MQTT configuration.
    mqttCfg MQTTConfig
    // The MQTT client.
    mqttClient mqtt.Client
}

// NewHQ returns new instance of HQ with some default values for fields.
func NewHQ() *HQ {
    return &HQ{
        detAgentIP:  "192.168.42.1",
        detUseIP:    "192.168.42.2",
        detCmdPort:  7802,
        detInterval: 5 * time.Second,
        cipher:      NewNoopCipher(),
        log:         NewLog(),
    }
}

// SetDet set WiFi interface name, agent's access point name regexp and password
// to use for detecting new agents.
func (hq *HQ) SetDet(itfName string, apNamePat *regexp.Regexp, apPass string) error {
    if itfName == "" {
        return errors.New("you must provide WiFi interface name")
    }
    hq.detItfName = itfName

    var err error
    hq.detItf, err = getInterface(hq.detItfName)
    if err != nil {
        return err
    }

    hq.detApNamePat = apNamePat
    hq.detApPass = apPass
    return nil
}

// SetDetIP set IP to use for WiFi interface after connecting to agent's
// access point (agents do not provide DHCP) and agent's IP.
func (hq *HQ) SetDetIPs(useIP, agentIP string) error {
    hq.detUseIP = useIP
    hq.detAgentIP = agentIP
    return nil
}

// SetDetCmdPort set TCP port to use for communication with agent.
// Every agent creates TC/IP command server during discovery phase to receive
// configuration for HQ access point and MQTT server.
func (hq *HQ) SetDetCmdPort(cmdPort int) error {
    hq.detCmdPort = cmdPort
    return nil
}

// SetDetInterval sets interval used for agent discovery retries.
func (hq *HQ) SetDetInterval(interval time.Duration) error {
    hq.detInterval = interval
    return nil
}

// SetCipher sets a cipher to use to encrypt and decrypt communication with agents.
func (hq *HQ) SetCipher(c Cipher) {
    hq.cipher = c
}

// SetAccessPoint sets HQ access point credentials which will be sent to
// new agent during discovery phase.
func (hq *HQ) SetAccessPoint(apName, apPass string) error {
    hq.mainApName = apName
    hq.mainApPass = apPass
    return nil
}

// SetMQTTBroker configures MQTT broker.
func (hq *HQ) SetMQTTClient(client mqtt.Client) error {
    hq.mqttClient = client
    return nil
}

// SetLogger set logger to use.
func (hq *HQ) SetLogger(l Logger) {
    hq.log = l
}

// Detect detects IoT access points in range.
func (hq *HQ) DetectAgents() ([]*AgentAP, error) {
    hq.log.Infof("scanning for new agents using %s interface...", hq.detItfName)
    aps, err := hq.detItf.scan()
    if err != nil {
        return nil, err
    }

    // Filter out non agent access points.
    var agents []*AgentAP
    for _, ap := range aps {
        if hq.detApNamePat.MatchString(ap.Name) {
            agents = append(agents, ap)
        }
    }

    if hq.isMQTTSet() {
        for _, agent := range agents {
            if err := hq.PublishMQTT("hq/new_agent", 0, false, agent.MAC()); err != nil {
                return agents, err
            }
        }
    }

    return agents, nil
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

    hq.mqttCfg = cfg
    hq.mqttClient = client
    return nil
}

// PublishMQTT publishes payload to MQTT broker.
func (hq *HQ) PublishMQTT(topic string, qos byte, retained bool, payload interface{}) error {
    hq.mqttClient.Publish(topic, qos, retained, payload).Wait()
    return nil
}

// GetConfigCmd returns configuration command.
func (hq *HQ) GetConfigCmd() *CmdConfig {
    cmd := NewConfigCmd()

    cmd.ApName = hq.mainApName
    cmd.ApPass = hq.mainApPass

    if hq.isMQTTSet() {
        cmd.MQTTIP = hq.mqttCfg.ip
        cmd.MQTTPort = hq.mqttCfg.port
        cmd.MQTTUser = hq.mqttCfg.user
        cmd.MQTTPass = hq.mqttCfg.pass
    }

    return cmd
}

// isMQTTSet returns true if MQTT configuration has been set.
func (hq *HQ) isMQTTSet() bool {
    return hq.mqttCfg.ip != "" && hq.mqttCfg.port != 0
}
