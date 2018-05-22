package cmd

import (
    "github.com/spf13/cobra"
    "github.com/pkg/errors"
)

var configureCmd = &cobra.Command{
    Use:   "configure",
    Short: "configure agent.",
    Long:  `configure agent.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        if len(args) != 1 {
            return errors.New("missing access point name argument")
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
