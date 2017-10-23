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

// +build linux

package iotdet

import (
    "github.com/pkg/errors"
    jww "github.com/spf13/jwalterweatherman"
    "io/ioutil"
    "net"
    "os"
    "os/exec"
    "regexp"
    "strings"
    "strconv"
)

// addressRegEx is regular expression for WiFi access point unique address.
var addressRegEx *regexp.Regexp = regexp.MustCompile("(:?.*?)Address: (([[:xdigit:]]{2}:){5}[[:xdigit:]]{2})")

// nameRegEx is regular expression for WiFi access point name.
var nameRegEx *regexp.Regexp = regexp.MustCompile("(:?.*?)ESSID:\"(.*?)\"")

func getWifiInterfaces() ([]net.Interface, error) {
    var err error
    var all []net.Interface
    var ret []net.Interface

    if all, err = net.Interfaces(); err != nil {
        return nil, errors.Wrap(err, "can't get interfaces")
    }

    for _, itf := range all {
        if dirExists("/sys/class/net/" + itf.Name + "/wireless") {
            jww.DEBUG.Println("Found " + itf.Name + " WiFi interface.")
            ret = append(ret, itf)
        }
    }

    return ret, nil
}

func scanForAPs(itf *wifiItf) ([]*AccessPoint, error) {
    var aps []*AccessPoint

    out, err := exec.Command("bash", "-c", "iwlist "+itf.Name+" scan | egrep 'ESSID:|Address:'").CombinedOutput()
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
        aps = append(aps, NewAccessPoint(name, mac, itf))
    }

    return aps, nil
}

// ifUp brings an interface up.
func ifUp(itf net.Interface) error {
    err := exec.Command("/sbin/ifconfig", itf.Name, "up").Run()
    if err != nil {
        if _, ok := err.(*exec.ExitError); ok {
            return err
        } else {
            return nil
        }
    }

    return nil
}

// isUp returns true if given WiFi interface is up.
func isUp(itf net.Interface) bool {
    _, err := exec.Command("bash", "-c", "/sbin/ifconfig | grep "+itf.Name).Output()
    if err != nil {
        if _, ok := err.(*exec.ExitError); ok {
            // The program has exited with an exit code != 0
            return false
        } else {
            return true
        }
    }

    return true
}

// connectToAp connects to given access point.
func connectToAp(apName string, apPass string, itfName string) (stopChanel, error) {
    var err error
    var wpaStopCh stopChanel

    stopCh := make(stopChanel)

    if err = killWpaDaemon(apName); err != nil {
        return nil, err
    }

    if err = configWrite(apName, apPass); err != nil {
        return nil, err
    }

    if wpaStopCh, err = startWpaDaemon(itfName); err != nil {
        return nil, err
    }

    go func() {
        select {
        case <-stopCh:
            wpaStopCh <- struct{}{}
            <-wpaStopCh
            close(stopCh)
            close(wpaStopCh)
        }
    }()

    return stopCh, nil
}

func setIp(itfName string, ip string) error {
    jww.DEBUG.Printf("Setting %s interface IP to %s.", itfName, ip)

    err := exec.Command("ifconfig", itfName, ip, "netmask", "255.255.255.0").Run()
    if err != nil {
        if _, ok := err.(*exec.ExitError); ok {
            return err
        } else {
            return nil
        }
    }

    return nil
}

func pingIot(itfName string, ip string) error {
    jww.DEBUG.Printf("Pinging IoT device at %s.", ip)

    err := exec.Command("ping", "-I", itfName, "-c1", "-W", "1", ip).Run()
    if err != nil {
        if _, ok := err.(*exec.ExitError); ok {
            return errors.Errorf("Can't reach %s with ping.", ip)
        } else {
            return nil
        }
    }

    return nil
}

// Checks if directory exists.
func dirExists(dirPath string) bool {
    if _, err := os.Stat(dirPath); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

// killWpaDaemon kills wpa_supplicant daemon.
func killWpaDaemon(itfName string) error {
    jww.DEBUG.Printf("Killing wpa_supplicant for %s interface.", itfName)
    exec.Command("bash", "-c", "wpa_cli -i "+itfName+" terminate").Run()

    return nil
}

// configWrite writes wpa_supplicant configuration file.
func configWrite(apName string, apPass string) error {
    cfgPath := getWpaConfigPath()
    jww.DEBUG.Printf("Writing wpa_supplicant daemon config to %s.\n", cfgPath)

    text := ""
    text += "ctrl_interface=/var/run/wpa_supplicant." + strconv.Itoa(os.Getpid()) + "\n"
    text += "network={" + "\n"
    text += "     ssid=\"" + apName + "\"" + "\n"
    text += "     psk=\"" + apPass + "\"\n"
    text += "}" + "\n"

    if err := ioutil.WriteFile(cfgPath, []byte(text), 0644); err != nil {
        return err
    }

    return nil
}

// getWpaConfigPath returns path to the wpa_supplicant configuration file.
func getWpaConfigPath() string {
    return os.TempDir() + "/iot_wpa_supplicant." + strconv.Itoa(os.Getpid()) + ".conf"
}

// startWpaDaemon starts wpa_supplicant daemon.
func startWpaDaemon(itfName string) (stopChanel, error) {
    jww.DEBUG.Println("Starting wpa_supplicant daemon.")

    // wpa_supplicant -i wlan0 -D nl80211 -c /tmp/iot_wpa_supplicant.123.conf
    itf := "-i" + itfName
    driver := "-D" + "nl80211"
    config := "-c" + getWpaConfigPath()

    cmd := exec.Command("wpa_supplicant", itf, driver, config)
    stopChan, err := runCmdBg(cmd)
    if err != nil {
        return nil, errors.Wrap(err, "wpa_supplicant")
    }

    return stopChan, nil
}
