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
    "github.com/sirupsen/logrus"
)

// Config represents HQ configuration.
type Config struct {
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
    Log      *logrus.Entry
    LogLevel int
}
