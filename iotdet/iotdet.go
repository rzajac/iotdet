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

package main

import (
    "github.com/codegangsta/cli"
    "github.com/pkg/errors"
    jww "github.com/spf13/jwalterweatherman"
    "os"
    "github.com/go-ini/ini"
    "encoding/hex"
    "github.com/rzajac/iotdet/pkg/iotdet"
)

func main() {
    // Configure CLI application.
    app := cli.NewApp()
    app.Name = "iotdet"
    app.Usage = "detect new IoT devices in range."
    app.Action = iotDetect
    app.Version = "0.0.1"
    app.Authors = []cli.Author{
        {Name: "Rafal Zajac", Email: "rzajac@gmail.com"},
    }

    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:  "interface, i",
            Usage: "wifi interface name",
        },
        cli.StringFlag{
            Name:  "config, c",
            Usage: "path to iotdet configuration file",
        },
        cli.IntFlag{
            Name:  "log level, l",
            Usage: "log level 0 - 6",
        },
    }

    if os.Getegid() != 0 {
        jww.ERROR.Println(app.Name + " needs root privileges.")
        os.Exit(1)
    }

    // Run application.
    app.Run(os.Args)
}

// iotDetect detects and configures new IoT devices.
func iotDetect(ctx *cli.Context) {
    var err error
    var aps []*iotdet.DevAP
    var cfg *iotdet.IotCfg

    jww.SetStdoutThreshold(jww.Threshold(ctx.GlobalInt("log level")))

    cfg, err = configure(ctx)
    if err != nil {
        jww.ERROR.Println(err)
        os.Exit(1)
    }

    if aps, err = iotdet.Detect(cfg); err != nil {
        jww.ERROR.Println(err)
        os.Exit(1)
    }

    if err = iotdet.Configure(cfg, aps); err != nil {
        jww.ERROR.Println(err)
        os.Exit(1)
    }

    jww.INFO.Println("Done.")
}

// configure loads iotdet configuration.
func configure(ctx *cli.Context) (*iotdet.IotCfg, error) {
    cfg := &iotdet.IotCfg{}
    itfName := ctx.GlobalString("interface")

    if itfName == "" {
        return nil, errors.New("you must provide WiFi interface name")
    }
    cfg.ItfName = itfName

    cfgFile, err := ini.Load(ctx.GlobalString("config"))
    if err != nil {
        return nil, err
    }

    cph := cfgFile.Section("iotdet").Key("cipher").In("none", []string{"aes"})
    if cph == "aes" {
        key, err := hex.DecodeString(cfgFile.Section("cipher_aes").Key("key").String())
        if err != nil {
            return nil, errors.Errorf("Invalid AES key value %s.")
        }

        vi, err := hex.DecodeString(cfgFile.Section("cipher_aes").Key("vi").String())
        if err != nil {
            return nil, errors.Errorf("Invalid AES vi value %s.")
        }
        cfg.Cipher = iotdet.NewAesDrv(key, vi)
    } else if cph == "none" {
        cfg.Cipher = &iotdet.Noop{}
    }
    cfg.IotIp = cfgFile.Section("iotdet").Key("iot_ip").String()
    cfg.MyIp = cfgFile.Section("iotdet").Key("my_ip").String()

    cfg.UdpPort = cfgFile.Section("udp").Key("port").MustInt()
    cfg.TcpPort = cfgFile.Section("cmd_server").Key("port").MustInt()
    cfg.IotApName = cfgFile.Section("iot_access_point").Key("name").String()
    cfg.IotApPass = cfgFile.Section("iot_access_point").Key("password").String()
    cfg.DevApPass = cfgFile.Section("iot_dev").Key("ap_password").String()

    return cfg, nil
}
