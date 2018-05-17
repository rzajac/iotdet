package cmd

import (
    "github.com/spf13/cobra"
    "github.com/pkg/errors"
)

var configureCmd = &cobra.Command{
    Use:   "configure",
    Short: "Configure new agent.",
    Long:  `Configure new agent.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        if len(args) != 1 {
            return errors.New("missing access point name argument")
        }

        //apName := args[0]

        return nil
    },
}

func init() {
    rootCmd.AddCommand(configureCmd)
}
