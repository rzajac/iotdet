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
    "strings"
)

// beacon represents access point agent creates during discovery phase.
type beacon struct {
    // Access point name which must have unique ID of the beacon as part of it.
    // It is constructed as two strings separated by underscore where the second
    // component is the unique ID. It's customary to use MAC address as the
    // unique ID.
    // The same ID must be used by the agent for all other communication.
    name string
    // Access point password.
    pass string
    // The IP agent assigns to itself when creating access point.
    // This need to be known because agents do not run DHCP service.
    ip string
    // The IP to set on the interface which is used to connect to access point.
    itfIP string
    // The TCP port agent listens on for configuration commands during
    // agent detection phase.
    cmdPort int
}

// newBeacon returns new beacon instance.
func newBeacon(name, pass, ip, itfIP string, cmdPort int) *beacon {
    return &beacon{
        name:    name,
        pass:    pass,
        ip:      ip,
        itfIP:   itfIP,
        cmdPort: cmdPort,
    }
}

// ID returns beacon unique ID.
func (ap *beacon) ID() string {
    return strings.Split(ap.name, "_")[1]
}
