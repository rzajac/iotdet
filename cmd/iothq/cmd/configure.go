package cmd

import (
    "github.com/spf13/cobra"
    "github.com/pkg/errors"
)

var configureCmd = &cobra.Command{
    Use:   "configure [interface name]",
    Short: "Configure agent.",
    Long:  `Configure agent.`,
    Example: "   configure wlp3s0\n   configure wlx000f55a93e30",
    RunE: func(cmd *cobra.Command, args []string) error {
        if len(args) != 1 {
            return errors.New("requires access point name argument")
        }

        h, err := getConfiguredHQ()
        if err != nil {
            return err
        }

        return h.Configure(args[0])
    },
}

func init() {
    rootCmd.AddCommand(configureCmd)
}
