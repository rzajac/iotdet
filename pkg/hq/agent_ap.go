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

// agentAP represents access point agent creates during discovery phase.
type agentAP struct {
    name  string // Access point name.
    pass  string // Access point password.
    ip    string // The IP agent assigns to itself when creating access point.
    port  int    // The TCP port agents listen on for configuration commands.
    useIP string // The IP to use after connecting to agent's access point.
}

// newAgentAP returns new agentAP instance.
func newAgentAP(name string) *agentAP {
    return &agentAP{
        name: name,
    }
}

// Mac returns access point mac address.
func (ap *agentAP) Mac() string {
    return strings.Split(ap.name, "_")[1]
}
