package cmd

import (
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/rzajac/iotdet/pkg/iotdet"
    "os"
    "time"
)

var startCmd = &cobra.Command{
    Use:   "start",
    Short: "Start IoT detection and configuration service.",
    Long:  `Start IoT detection and configuration service.`,
    Run: func(cmd *cobra.Command, args []string) {
        itfName := viper.GetString("iotdet.itf_name")

        ctrl := make(iotdet.CtrlChanel)

        go func() {
            for {
                <-time.After(1 * time.Second)
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

        detector, err := iotdet.NewDetector(itfName, log)
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
