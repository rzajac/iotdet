package cmd

import (
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var startCmd = &cobra.Command{
    Use:   "start",
    Short: "Start IoT detection and configuration service.",
    Long:  `Start IoT detection and configuration service.`,
    Run: func(cmd *cobra.Command, args []string) {
        host := viper.GetString("auth.db.auth.host")
    },
}

func init() {
    rootCmd.AddCommand(startCmd)
}
