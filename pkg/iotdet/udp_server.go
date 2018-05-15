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
    "github.com/pkg/errors"
    "net"
    "strconv"
    "github.com/sirupsen/logrus"
)

// UdpServer represents UDP server listening for
// discovery broadcasts from IoT devices.
type UdpServer struct {
    cfg  *IotCfg
    itf  *net.Interface
    stop chan bool
    log  *logrus.Entry
}

func NewUdpServer(cfg *IotCfg, itf *net.Interface, log *logrus.Entry) *UdpServer {
    return &UdpServer{
        cfg: cfg,
        itf: itf,
        log: log,
    }
}

// Start starts UDP server listening for discovery broadcasts on given port.
func (udp *UdpServer) Start(stop chan bool) error {
    var err error
    var srvAddr *net.UDPAddr
    var srvConn *net.UDPConn

    if udp.stop != nil {
        return errors.New("UDP server already started")
    }

    udp.stop = stop

    // Setup UDP address.
    udp.log.Infof("starting UDP server on port %d", strconv.Itoa(udp.cfg.UdpPort))
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
                udp.log.Infof("stopping UDP server")
                return

            default:
                n, addr, err := srvConn.ReadFromUDP(buf)
                if err != nil {
                    udp.log.Error("UDP: ", err)
                    continue
                }

                b := newIotDev(addr.IP.String(), &Noop{})
                err = b.parseDiscoveryBroadcast(buf[0:n])
                if err != nil {
                    udp.log.Error(err)
                }

                _, err = udp.handleDiscovery(b)
                if err != nil {
                    udp.log.Error(err)
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
        return nil, errors.Errorf("%s has no IP addresses", udp.itf.Name)
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
