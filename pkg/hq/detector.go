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
    "time"
)

// Detector is responsible for detecting and configuring new IoT devices.
type Detector struct {
    cfg  *HQ
    itf  *WiFiItf
    ctrl CtrlChanel
}

// NewDetector returns new Detector instance.
func NewDetector(cfg *HQ) (*Detector, error) {
    itf, err := GetInterface(cfg)
    if err != nil {
        return nil, err
    }
    return &Detector{
        cfg:  cfg,
        itf:  itf,
        ctrl: make(CtrlChanel)}, nil
}

// Detect detects IoT access points in range.
func (d *Detector) Detect() ([]*AgentAP, error) {
    d.cfg.Log.Infof("Scanning for new IoT devices using %s interface...", d.itf.Name())
    return d.itf.Scan()
}

// Start starts detection service.
func (d *Detector) Start() error {
    go func() {
        for {
            <-time.After(d.cfg.DetInterval)
            select {
            case cmd := <-d.ctrl:
                if cmd == "STOP" {
                    d.cfg.Log.Info("stopping detector")
                    return
                }

            default:
                aps, err := d.Detect()
                if err != nil {
                    d.cfg.Log.Error(err)
                    continue
                }
                // TODO: configure agents.
                d.cfg.Log.Infof("found %d agents", len(aps))
            }
        }
    }()

    return nil
}

// Stop stops detection service.
func (d *Detector) Stop() {
    d.ctrl <- "STOP"
}

// Detect is a helper function which detects new agents.
func Detect(cfg *HQ) ([]*AgentAP, error) {
    detector, err := NewDetector(cfg)
    if err != nil {
        return nil, err
    }

    return detector.Detect()
}