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

// Represents MQTT broker configuration.
type MQTTConfig struct {
    // MQTT client ID.
    clientID string
    // MQTT broker IP.
    ip string
    // MQTT broker port.
    port int
    // MQTT broker username.
    user string
    // MQTT broker password.
    pass string
}

// NewMQTTConfig returns new MQTT broker configuration struct.
func NewMQTTConfig(clientID, ip string, port int, user, pass string) (MQTTConfig, error) {
    return MQTTConfig{
        clientID: clientID,
        ip:       ip,
        port:     port,
        user:     user,
        pass:     pass,
    }, nil
}
