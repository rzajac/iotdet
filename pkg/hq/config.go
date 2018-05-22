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

import (
    "regexp"
    "time"
)

// Configurator is an interface used to configure HQ.
type Configurator interface {
    // GetDetItfName returns the network interface name to
    // use for beacon detection (new agents).
    GetDetItfName() string
    // GetItfIP returns the IP to use after connecting to a beacon.
    GetItfIP() string
    // GetDetInterval returns the interval for scanning for new beacons.
    GetDetInterval() time.Duration
    // GetBeaconNamePat returns the regexp used to match beacon names.
    GetBeaconNamePat() *regexp.Regexp
    // GetBeaconPass returns the password to use when connecting to beacons.
    GetBeaconPass() string
    // GetBeaconIP returns the IP beacon assigns to itself.
    GetBeaconIP() string
    // GetBeaconPort returns the TCP port beacons listen on for
    // configuration commands.
    GetBeaconPort() int
    // GetCipher returns the cipher to use for communication with beacons.
    GetCipher() Cipher
    // GetAPName returns the access point name agents use to communicate.
    GetAPName() string
    // GetAPPass returns the access point password.
    GetAPPass() string
    // GetMQTTIP returns MQTT broker IP.
    GetMQTTIP() string
    // GetMQTTPort returns MQTT broker port.
    GetMQTTPort() int
    // GetMQTTUser returns  MQTT broker username.
    GetMQTTUser() string
    // GetMQTTPass returns MQTT broker password.
    GetMQTTPass() string
    // GetMQTTClientID returns MQTT client ID to use in HQ.
    GetMQTTClientID() string
}
