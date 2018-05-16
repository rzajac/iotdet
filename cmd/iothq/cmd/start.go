package cmd

import (
    "github.com/spf13/cobra"
    "os"
    "time"
    "github.com/rzajac/iotdet/pkg/hq"
)

var startCmd = &cobra.Command{
    Use:   "start",
    Short: "Start IoT detection and configuration service.",
    Long:  `Start IoT detection and configuration service.`,
    Run: func(cmd *cobra.Command, args []string) {
        cfg, err := config()
        if err != nil {
            log.Error(err)
            os.Exit(1)
        }

        ctrl := make(hq.CtrlChanel)

        go func() {
            for {
                <-time.After(cfg.DetInterval)
                select {
                case cmd := <-ctrl:
                    if cmd == "STOP" {
                        log.Info("STOPPING")
                        return
                    }

                default:
                    log.Info("DEFAULT")
                }
            }
        }()

        <-time.After(10 * time.Second)
        ctrl <- "STOP"

        detector, err := hq.NewDetector(cfg)
        if err != nil {
            log.Error(err)
            os.Exit(1)
        }

        if _, err := detector.Detect(); err != nil {
            log.Error(err)
            os.Exit(1)
        }

        //if err := iotdet.Configure(cfg, aps); err != nil {
        //    logrus.Error(err)
        //    os.Exit(1)
        //}
    },
}

func init() {
    rootCmd.AddCommand(startCmd)
}
