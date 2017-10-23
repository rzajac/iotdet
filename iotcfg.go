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

package iotdet

// IotCfg represents IoT detection configuration.
type IotCfg struct {
    ItfName   string // Represents network interface to scan for new IoT devices.
    MyIp      string // The IP to assign to wifi interface.
    LogLevel  int    // Log level.
    Cipher    Cipher // The cipher to use for communication with IoT devices.
    UdpPort   int    // The UDP port IoT device uses for Main Server detection broadcasts.
    TcpPort   int    // The TCP port IoT device listens for commands on.
    IotApName string // The access point for IoT devices.
    IotApPass string // The IotApName password.
    IotIp     string // The IoT device IP.
    DevApPass string // The IoT device access point password.
    AesKey    []byte // AES encryption key
    AesVi     []byte // AES CBC initialization vector.
}
