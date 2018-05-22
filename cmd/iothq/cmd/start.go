package cmd

import (
    "github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
    Use:   "start",
    Short: "Start HomeHQ service",
    Long:  `Start HomeHQ service.`,
    Hidden: true,
    RunE: func(cmd *cobra.Command, args []string) error {
        // Handle CTRL+C.
        //sig := make(chan os.Signal, 2)
        //signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
        //
        //det, err := hq.NewDetector(cfg)
        //if err != nil {
        //    return err
        //}
        //
        //// Start new agent detection goroutine.
        //if err := det.Start(); err != nil {
        //    return err
        //}
        //
        //<-sig
        //det.Stop()

        return nil
    },
}

func init() {
    rootCmd.AddCommand(startCmd)
}
