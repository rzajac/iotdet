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
    "encoding/json"
    "net"
    "time"
    "bytes"
    "io"
    "strconv"
)

// Agent represents an agent.
type Agent struct {
    Name    string // The agent's name.
    IP      string // The agent's IP address.
    CmdPort int    // The agent's command server port.

    cfg *Config
}

// NewAgent returns new instance of Agent.
func NewAgent(name, ip string, cfg *Config) *Agent {
    return &Agent{
        Name: name,
        IP:   ip,
        cfg:  cfg,
    }
}

// SendCmd sends command to the agent.
func (a *Agent) SendCmd(cmd Cmd) ([]byte, error) {
    c, err := json.Marshal(cmd)
    if err != nil {
        return nil, err
    }
    a.cfg.Log.Debugf("sending to %s command: %s", a.Name, string(c))

    return a.SendMsg(c)
}

// SendMsg sends message to the agent.
func (a *Agent) SendMsg(msg []byte) ([]byte, error) {
    var err error
    var resp []byte
    var conn net.Conn

    // connect to TCP server.
    if conn, err = a.connect(); err != nil {
        return resp, err
    }
    defer conn.Close()
    conn.SetReadDeadline(time.Now().Add(3 * time.Second))

    msg, err = a.cfg.Cipher.Encrypt(msg)
    if err != nil {
        return resp, err
    }

    _, err = conn.Write(msg)
    if err != nil {
        return resp, err
    }

    // Get the response.
    var buff bytes.Buffer
    _, err = io.Copy(&buff, conn)
    if err != nil {
        return resp, err
    }

    // Decrypt the response.
    resp, err = a.cfg.Cipher.Decrypt(buff.Bytes())
    if err != nil {
        return resp, err
    }

    a.cfg.Log.Debugf("agent responded with: %s", string(resp))

    return resp, err
}

// connect establishes TCP connection to agent.
func (a *Agent) connect() (net.Conn, error) {
    var err error
    var conn net.Conn

    // Build TCP server address.
    address := a.IP + ":" + strconv.Itoa(a.CmdPort)
    a.cfg.Log.Debugf("dialing agent %s" + address)

    // connect to TCP server.
    if conn, err = net.Dial("tcp", address); err != nil {
        return nil, err
    }

    return conn, err
}

// TODO
//func (iot *iotDev) parseDiscoveryBroadcast(json_msg []byte) error {
//    var disc discoveryCmd
//
//    jww.DEBUG.Printf("UDP message: %s\nFrom: %s", string(json_msg), iot.Ip)
//    if err := json.Unmarshal(json_msg, &disc); err != nil {
//        return err
//    }
//
//    if disc.Cmd != "iotDiscovery" {
//        return errors.New("not discovery broadcast")
//    }
//
//    return nil
//}