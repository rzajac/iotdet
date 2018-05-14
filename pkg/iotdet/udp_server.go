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
    "fmt"
    "github.com/pkg/errors"
    jww "github.com/spf13/jwalterweatherman"
    "net"
    "strconv"
)

// UdpServer represents UDP server listening for
// discovery broadcasts from IoT devices.
type UdpServer struct {
    cfg  *IotCfg
    itf  *net.Interface
    Stop chan bool
}

func NewUdpServer(cfg *IotCfg, itf *net.Interface) *UdpServer {
    return &UdpServer{
        cfg: cfg,
        itf: itf,
    }
}

// Start starts UDP server listening for discovery broadcasts on given port.
func (udp *UdpServer) Start(stop chan bool) error {
    var err error
    var srvAddr *net.UDPAddr
    var srvConn *net.UDPConn

    if udp.Stop != nil {
        return errors.New("UDP server already started.")
    }

    udp.Stop = stop

    // Setup UDP address.
    jww.FEEDBACK.Println("Starting UDP server " + ":" + strconv.Itoa(udp.cfg.UdpPort))
    if srvAddr, err = net.ResolveUDPAddr("udp", ":"+strconv.Itoa(udp.cfg.UdpPort)); err != nil {
        return err
    }

    // Start UDP server.
    if srvConn, err = net.ListenUDP("udp", srvAddr); err != nil {
        return err
    }

    go func() {
        buf := make([]byte, JSON_MAX_LENGTH)
        defer srvConn.Close()

        for {
            select {
            case <-stop:
                jww.FEEDBACK.Println("Stopping UDP server.")
                return

            default:
                n, addr, err := srvConn.ReadFromUDP(buf)
                if err != nil {
                    jww.ERROR.Println("UDP Server error: ", err)
                    continue
                }

                b := newIotDev(addr.IP.String(), &Noop{})
                err = b.parseDiscoveryBroadcast(buf[0:n])
                if err != nil {
                    jww.ERROR.Println(err)
                }

                _, err = udp.handleDiscovery(b)
                if err != nil {
                    jww.ERROR.Println(err)
                }
            }
        }
    }()

    return nil
}

// handleDiscovery handles discovery broadcast.
func (udp *UdpServer) handleDiscovery(dev *iotDev) ([]byte, error) {
    var err error
    var ips []net.IP

    ips, err = udp.getItfIps(udp.itf)
    if err != nil {
        return nil, err
    }

    if len(ips) == 0 {
        msg := fmt.Sprintf("WiFiItf %f has no IP addresses.\n", udp.itf.Name)
        return nil, errors.New(msg)
    }

    return dev.sendCmd(udp.cfg.UdpPort, newServerCmd(ips[0], udp.cfg.TcpPort))
}

// getItfIps returns given interface IPs.
func (udp *UdpServer) getItfIps(itf *net.Interface) ([]net.IP, error) {
    var ip net.IP
    var ips []net.IP

    addresses, err := itf.Addrs()
    if err != nil {
        return ips, err
    }

    for _, address := range addresses {
        switch v := address.(type) {
        case *net.IPNet:
            ip = v.IP
        case *net.IPAddr:
            ip = v.IP
        }
    }

    ip = ip.To4()
    if ip != nil {
        ips = append(ips, ip)
    }

    return ips, nil
}
