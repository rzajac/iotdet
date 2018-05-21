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

import "encoding/json"

// cmdConfig represents configuration structure agents
// expect to receive during detection phase.
type cmdConfig struct {
    // The command identifier.
    Cmd      string `json:"cmd"`
    // The access point name the agent should use for communication.
    ApName   string `json:"ap_name"`
    // The access point password.
    ApPass   string `json:"ap_pass"`
    // The MQTT broker address.
    MQTTIP   string `json:"mqtt_ip"`
    // The MQTT broker port.
    MQTTPort int    `json:"mqtt_port"`
    // The MQTT broker username.
    MQTTUser string `json:"mqtt_user"`
    // The MQTT broker password.
    MQTTPass string `json:"mqtt_pass"`
}

// newConfigCmd returns configuration command.
func newConfigCmd() *cmdConfig {
    return &cmdConfig{Cmd: "cfg"}
}

func (c *cmdConfig) MarshalCmd() []byte {
    cmd, _ := json.Marshal(c)
    return cmd
}
