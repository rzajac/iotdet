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

import (
    "encoding/json"
    "github.com/pkg/errors"
    jww "github.com/spf13/jwalterweatherman"
    "net"
    "strconv"
    "time"
    "bytes"
    "io"
)

// IoT device representation.
type iotDev struct {
    Ip    string
    Crypt Cipher
}

func newIotDev(ip string, c Cipher) *iotDev {
    return &iotDev{
        Ip:    ip,
        Crypt: c,
    }
}

func (iot *iotDev) tcpConnect(port int) (net.Conn, error) {
    var err error
    var conn net.Conn

    // Build TCP server address.
    address := iot.Ip + ":" + strconv.Itoa(port)
    jww.DEBUG.Println("Dialing " + address)

    // connect to TCP server.
    if conn, err = net.Dial("tcp", address); err != nil {
        jww.ERROR.Println(err)
        return nil, err
    }

    return conn, err
}

func (iot *iotDev) sendMsg(port int, msg []byte) ([]byte, error) {
    var err error
    var resp []byte
    var conn net.Conn

    // connect to TCP server.
    if conn, err = iot.tcpConnect(port); err != nil {
        return resp, err
    }
    defer conn.Close()
    conn.SetReadDeadline(time.Now().Add(3 * time.Second))

    //jww.DEBUG.Println("Sending:")
    //dumpBytes(msg)

    msg, err = iot.Crypt.Encrypt(msg)
    if err != nil {
        return resp, err
    }

    //jww.DEBUG.Println("Sending enc:")
    //dumpBytes(msg)

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

    //jww.DEBUG.Println("Received:")
    //dumpBytes(buff.Bytes())

    // Decrypt the response.
    resp, err = iot.Crypt.Decrypt(buff.Bytes())
    if err != nil {
        return resp, err
    }

    //jww.DEBUG.Println("Decoded:")
    //dumpBytes(resp)

    jww.FEEDBACK.Printf("iotDev responded with: " + string(resp))

    return resp, err
}

func (iot *iotDev) sendCmd(port int, cmd interface{}) ([]byte, error) {
    c, err := json.Marshal(cmd)
    if err != nil {
        return nil, err
    }

    jww.DEBUG.Println("Sending command:", string(c))

    return iot.sendMsg(port, c)
}

func (iot *iotDev) parseDiscoveryBroadcast(json_msg []byte) error {
    var disc discoveryCmd

    jww.DEBUG.Printf("UDP message: %s\nFrom: %s", string(json_msg), iot.Ip)
    if err := json.Unmarshal(json_msg, &disc); err != nil {
        return err
    }

    if disc.Cmd != "iotDiscovery" {
        return errors.New("not discovery broadcast")
    }

    return nil
}
