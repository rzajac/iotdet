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
    "fmt"
    "time"
    "github.com/eclipse/paho.mqtt.golang"
)

func NewMQTTClient(cfg *HQ) (mqtt.Client, error) {
    opts := mqtt.NewClientOptions()
    opts.AddBroker(fmt.Sprintf("tcp://%s:%d", cfg.MQTTIP, cfg.MQTTPort))
    opts.SetClientID("iothq")
    opts.SetKeepAlive(2 * time.Second)
    opts.SetPingTimeout(1 * time.Second)
    opts.SetUsername(cfg.MQTTUser)
    opts.SetPassword(cfg.MQTTPass)

    c := mqtt.NewClient(opts)
    if token := c.Connect(); token.Wait() && token.Error() != nil {
        return nil, token.Error()
    }
    return c, nil
}
