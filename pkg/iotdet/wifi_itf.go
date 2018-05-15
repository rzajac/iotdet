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
    "time"
    "sync"
    "github.com/sirupsen/logrus"
    "os/exec"
    "strconv"
    "os"
    "io/ioutil"
    "bufio"
    "strings"
    "regexp"
)

// addressRegEx is regular expression for WiFi access point unique address.
var addressRegEx = regexp.MustCompile("(:?.*?)Address: (([[:xdigit:]]{2}:){5}[[:xdigit:]]{2})")

// nameRegEx is regular expression for WiFi access point name.
var nameRegEx = regexp.MustCompile("(:?.*?)ESSID:\"(.*?)\"")

// interfaces is a collection of WiFi interfaces.
var interfaces map[string]*WiFiItf

func init() {
    interfaces = make(map[string]*WiFiItf, 2)
}

// GetInterface returns WiFi interface or error if it does not exist.
func GetInterface(name string, log *logrus.Entry) (*WiFiItf, error) {
    if itf, ok := interfaces[name]; ok {
        return itf, nil
    }

    // Check if given interface name exists on the system.
    if !dirExists("/sys/class/net/" + name + "/wireless") {
        return nil, errors.Errorf("wifi interface %s not found", name)
    }

    var err error
    var itf []net.Interface
    if itf, err = net.Interfaces(); err != nil {
        return nil, errors.Wrap(err, "can't get interfaces")
    }

    for _, itf := range itf {
        if itf.Name == name {
            interfaces[name] = &WiFiItf{
                Mutex:  &sync.Mutex{},
                itf:    itf,
                log:    log,
                discCh: make(stopChanel),
            }
            return interfaces[name], nil
        }
    }

    return nil, errors.Errorf("wifi interface %s not found", name)
}

// WiFiItf represents WiFi interface.
type WiFiItf struct {
    *sync.Mutex
    itf    net.Interface
    log    *logrus.Entry
    discCh stopChanel
}

// Configure configures IoT devices.
func (w *WiFiItf) Configure(myIP, iotIP string, iotPort int, name, pass string, cipher *Cipher) error {
    var err error
    var aps []*DevAP

    if aps, err = w.scan(); err != nil {
        return err
    }

    var iotDev *iotDev
    for _, ap := range aps {
        if err = ap.Connect(pass); err != nil {
            w.log.Error(err)
            ap.Disconnect()
            continue
        }

        if err = ap.Itf.setIP(myIP); err != nil {
            ap.Disconnect()
            return err
        }

        if err = ap.Itf.ping(iotIP); err != nil {
            ap.Disconnect()
            return err
        }

        iotDev = newIotDev(iotIP, cipher)
        if _, err = iotDev.sendCmd(iotPort, newApCmd(name, pass)); err != nil {
            w.log.Error(err)
        }

        ap.Disconnect()
    }

    return nil
}

// setIP sets IP on the interface.
func (w *WiFiItf) setIP(ip string) error {
    w.log.Debugf("setting %s IP to %s", w.itf.Name, ip)
    err := exec.Command("ifconfig", w.itf.Name, ip, "netmask", "255.255.255.0").Run()
    if err != nil {
        if _, ok := err.(*exec.ExitError); ok {
            return err
        }
    }

    return nil
}

// ping pings IP address and returns error if IP cannot be pinged.
func (w *WiFiItf) ping(ip string) error {
    w.log.Debugf("pinging device at %s", ip)
    err := exec.Command("ping", "-I", w.itf.Name, "-c1", "-W", "1", ip).Run()
    if err != nil {
        if _, ok := err.(*exec.ExitError); ok {
            return errors.Errorf("cannot reach %s with ping", ip)
        }
    }

    return nil
}

// connect connects to access point.
func (w *WiFiItf) connect(apName, apPass string) error {
    w.log.Debugf("killing wpa_supplicant for %s interface.", w.itf.Name)
    exec.Command("bash", "-c", "wpa_cli -i "+w.itf.Name+" terminate").Run()
    if err := w.wpaConfigWrite(apName, apPass); err != nil {
        return err
    }

    var wpaStopCh stopChanel
    var wpaConnCh connChanel
    var err error
    if wpaConnCh, wpaStopCh, err = w.wpaStartDaemon(); err != nil {
        return err
    }

    go func() {
        select {
        case <-w.discCh:
            wpaStopCh <- struct{}{}
            <-wpaStopCh
            close(wpaStopCh)
            close(w.discCh)
        }
    }()

    w.log.Debugf("connecting to %s with %s", apName, w.itf.Name)
    select {
    case <-wpaConnCh:
        w.log.Info("connected to %s with %s", apName, w.itf.Name)
        break
    case <-time.After(10 * time.Second):
        w.discCh <- struct{}{}
        <-w.discCh
        return errors.Errorf("connection to %s with %s timed out", apName, w.itf.Name)
    }

    return nil
}

// disconnect disconnects from access point.
func (w *WiFiItf) disconnect() {
    select {
    case w.discCh <- struct{}{}:
        w.log.Debugf("disconnecting %s", w.itf.Name)
        <-w.discCh
    default:
        return
    }
}

// isUp returns true if interface is up.
func (w *WiFiItf) isUp() bool {
    _, err := exec.Command("bash", "-c", "/sbin/ifconfig | grep "+w.itf.Name).Output()
    if err != nil {
        if _, ok := err.(*exec.ExitError); ok {
            // The program has exited with an exit code != 0
            return false
        }
    }

    return true
}

// up brings the interface up.
func (w *WiFiItf) up() error {
    err := exec.Command("/sbin/ifconfig", w.itf.Name, "up").Run()
    if err != nil {
        if _, ok := err.(*exec.ExitError); ok {
            return err
        }
    }

    return nil
}

// makeSureIsUp makes sure the interface is up.
func (w *WiFiItf) makeSureIsUp() error {
    if w.isUp() {
        return nil
    }

    w.log.Debugf("waiting for %s to became available", w.itf.Name)
    if err := w.up(); err != nil {
        return err
    }

    stopCh := runUntil(func() bool {
        return w.isUp()
    }, 1*time.Second, 5)

    success := <-stopCh
    if !success {
        return errors.Errorf("could not bring up %s", w.itf.Name)
    }

    w.log.Debugf("% is up", w.itf.Name)

    return nil
}

// scan returns a list of IoT WiFi access points in range.
func (w *WiFiItf) scan() ([]*DevAP, error) {
    var aps []*DevAP

    if err := w.makeSureIsUp(); err != nil {
        return aps, err
    }

    out, err := exec.Command("bash", "-c", "iwlist "+w.itf.Name+" scan | egrep 'ESSID:|Address:'").CombinedOutput()
    if err != nil {
        return aps, err
    }

    lines := strings.Split(string(out), "\n")
    lineCount := len(lines)

    for i := 0; i < lineCount; i += 2 {
        if i+1 > lineCount-1 {
            break
        }

        mac := strings.TrimSpace(addressRegEx.FindAllStringSubmatch(lines[i], -1)[0][2])
        name := strings.TrimSpace(nameRegEx.FindAllStringSubmatch(lines[i+1], -1)[0][2])

        w.log.Debugf("found %s access point", name)
        dev := NewDevAP(name, mac, w)

        if dev.IsIotAp() {
            aps = append(aps, NewDevAP(name, mac, w))
        }
    }

    return aps, nil
}

// wpaConfigWrite writes wpa_supplicant configuration file.
func (w *WiFiItf) wpaConfigWrite(apName string, apPass string) error {
    w.log.Debug("writing wpa_supplicant daemon config to %s", w.wpaConfigPath())

    text := ""
    text += "ctrl_interface=/var/run/wpa_supplicant." + strconv.Itoa(os.Getpid()) + "\n"
    text += "network={" + "\n"
    text += "     ssid=\"" + apName + "\"" + "\n"
    text += "     psk=\"" + apPass + "\"\n"
    text += "}" + "\n"

    return ioutil.WriteFile(w.wpaConfigPath(), []byte(text), 0644)
}

// wpaConfigPath returns path to the wpa_supplicant configuration file.
func (w *WiFiItf) wpaConfigPath() string {
    return os.TempDir() + "/iot_wpa_supplicant." + strconv.Itoa(os.Getpid()) + ".conf"
}

// startWpaDaemon starts wpa_supplicant daemon.
func (w *WiFiItf) wpaStartDaemon() (connChanel, stopChanel, error) {
    w.log.Debug("starting wpa_supplicant daemon")

    // wpa_supplicant -i wlan0 -D nl80211 -c /tmp/iot_wpa_supplicant.123.conf
    itf := "-i" + w.itf.Name
    driver := "-D" + "nl80211"
    config := "-c" + w.wpaConfigPath()

    cmd := exec.Command("wpa_supplicant", itf, driver, config)

    stopCh := make(stopChanel)
    connCh := make(connChanel)
    stdoutCh := make(chan string)

    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return nil, nil, err
    }

    if err := cmd.Start(); err != nil {
        return nil, nil, err
    }

    go func() {
        scanner := bufio.NewScanner(stdout)
        for scanner.Scan() {
            stdoutCh <- scanner.Text()
        }
    }()

    go func() {
        for {
            select {
            case t := <-stdoutCh:
                w.log.Debugf("WPA: %s", t)
                if strings.Contains(t, "CTRL-EVENT-CONNECTED") {
                    connCh <- struct{}{}
                }
            case <-stopCh:
                w.log.Debugf("killing PID %d.", cmd.Process.Pid)
                cmd.Process.Kill()
                path := w.wpaConfigPath()
                w.log.Debugf("removing wpa_supplicant daemon config file %s", path)
                os.Remove(path)
                stopCh <- struct{}{}
                return
            }
        }
    }()

    return connCh, stopCh, nil
}

// wpaKillDaemon kills wpa_supplicant daemon.
func (w *WiFiItf) wpaKillDaemon() error {
    w.log.Debug("killing wpa_supplicant for %s interface", w.itf.Name)
    exec.Command("bash", "-c", "wpa_cli -i "+w.itf.Name+" terminate").Run()

    return nil
}
