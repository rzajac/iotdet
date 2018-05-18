package cmd

import (
    "github.com/spf13/cobra"
    "github.com/rzajac/iotdet/pkg/hq"
    "os"
    "os/signal"
    "syscall"
)

var startCmd = &cobra.Command{
    Use:   "start",
    Short: "Start IoT detection and configuration service.",
    Long:  `Start IoT detection and configuration service.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        // Handle CTRL+C.
        sig := make(chan os.Signal, 2)
        signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

        det, err := hq.NewDetector(cfg)
        if err != nil {
            return err
        }

        // Start new agent detection goroutine.
        if err := det.Start(); err != nil {
            return err
        }

        <-sig
        det.Stop()

        return nil
    },
}

func init() {
    rootCmd.AddCommand(startCmd)
}
