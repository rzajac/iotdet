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

// AgentAP represents access point agent creates during discovery phase.
type AgentAP struct {
    Name string // Access point name.
}

// NewAgentAP returns new AgentAP instance.
func NewAgentAP(name string) *AgentAP {
    return &AgentAP{Name: name}
}

// MAC returns access point MAC address.
func (ap *AgentAP) MAC() string {
    return strings.Split(ap.Name, "_")[1]
}
