package cmd

import (
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/rzajac/iotdet/pkg/iotdet"
    "os"
)

var startCmd = &cobra.Command{
    Use:   "start",
    Short: "Start IoT detection and configuration service.",
    Long:  `Start IoT detection and configuration service.`,
    Run: func(cmd *cobra.Command, args []string) {
        itfName := viper.GetString("iotdet.itf_name")

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
