package cmd

import (
    "github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
    Use:   "start",
    Short: "Start IoT detection and configuration service.",
    Long:  `Start IoT detection and configuration service.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        _, err := config()
        if err != nil {
            return err
        }
        return nil
    },
}

func init() {
    rootCmd.AddCommand(startCmd)
}
