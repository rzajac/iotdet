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
)

// HQ represents HQ configuration.
type HQ struct {
    DetItfName  string        // The interface name to use for new agents detection.
    DetApPass   string        // The password for agent's access point.
    DetAgentIP  string        // The IP agent assigns to itself when creating access point.
    DetUseIP    string        // The IP to use after connecting to agent's access point.
    DetCmdPort  int           // The TCP port agents listen on for configuration commands.
    DetInterval time.Duration // The interval for scanning for new agents.

    Cipher Cipher // The cipher to use for communication with agent devices.

    HQApName string // The HQ access point name agents use to communicate.
    HQApPass string // The HQ access point password.

    // MQTT broker configuration.
    MQTTIP   string
    MQTTPort int
    MQTTUser string
    MQTTPass string

    // Logger configuration.
    Log Logger
}

// NewHQ returns new instance of HQ with some default values for fields.
func NewHQ() *HQ {
    return &HQ{
        DetAgentIP:  "192.168.42.1",
        DetUseIP:    "192.168.42.2",
        DetCmdPort:  7802,
        DetInterval: 5 * time.Second,
        Cipher:      NewNoopCipher(),
        Log:         NewLog(),
    }
}

// SetDetPass set wifi interface name and access point password to use for
// detecting new agents.
func (hq *HQ) SetDetPass(itfName, apPass string) {
    hq.DetItfName = itfName
    hq.DetApPass = apPass

}

// SetDetPass set IP to use for WiFi interface after connecting to agent's
// access point (agents do not provide DHCP) and agent's IP.
func (hq *HQ) SetDetIPs(useIP, agentIP string) {
    hq.DetUseIP = useIP
    hq.DetAgentIP = agentIP
}

// SetDetCmdPort set TCP port to use for communication with agent.
// Every agent creates TC/IP command server during discovery phase to receive
// configuration for HQ access point and MQTT server.
func (hq *HQ) SetDetCmdPort(cmdPort int) {
    hq.DetCmdPort = cmdPort
}

// SetDetInterval sets interval used for agent discovery retries.
func (hq *HQ) SetDetInterval(interval time.Duration) {
    hq.DetInterval = interval
}

// SetCipher sets a cipher to use to encrypt and decrypt communication with agents.
func (hq *HQ) SetCipher(c Cipher) {
    hq.Cipher = c
}

// SetHQAccessPoint sets HQ access point credentials which will be sent to
// new agent during discovery phase.
func (hq *HQ) SetHQAccessPoint(apName, apPass string) {
    hq.HQApName = apName
    hq.HQApPass = apPass
}

// SetMQTTBroker configures MQTT broker.
func (hq *HQ) SetMQTTBroker(ip string, port int, user, pass string) {
    hq.MQTTIP = ip
    hq.MQTTPort = port
    hq.MQTTUser = user
    hq.MQTTPass = pass
}

// GetConfigCmd returns configuration command.
func (hq *HQ) GetConfigCmd() *CmdConfig {
    cmd := NewConfigCmd()

    cmd.ApName = hq.HQApName
    cmd.ApPass = hq.HQApPass
    cmd.MQTTIP = hq.MQTTIP
    cmd.MQTTPort = hq.MQTTPort
    cmd.MQTTUser = hq.MQTTUser
    cmd.MQTTPass = hq.MQTTPass

    return cmd
}
