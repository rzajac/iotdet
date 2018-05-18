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

// The maximum expected command length.
const CMD_MAX_LENGTH = 512

// MarshalCmd is an interface which all commands must implement.
type MarshalCmd interface {
    // MarshalCmd marshals command.
    MarshalCmd() []byte
}

// MarshalCmd is an interface which all commands must implement.
type UnmarshalCmd interface {
    // MarshallCmd marshals command.
    UnmarshalCmd(data []byte)
}
