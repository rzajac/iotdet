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

// CmdConfig represents structure agents
// expect to receive during detection phase.
type CmdConfig struct {
    Cmd      string `json:"cmd"`

    ApName   string `json:"ap_name"`
    ApPass   string `json:"ap_pass"`

    MQTTIP   string `json:"mqtt_ip"`
    MQTTPort int    `json:"mqtt_port"`
    MQTTUser string `json:"mqtt_user"`
    MQTTPass string `json:"mqtt_pass"`
}

// NewConfigCmd returns configuration command.
func NewConfigCmd() *CmdConfig {
    return &CmdConfig{Cmd: "cfg"}
}

func (c *CmdConfig) MarshalCmd() []byte {
    cmd, _ := json.Marshal(c)
    return cmd
}
