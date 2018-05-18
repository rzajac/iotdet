package cmd

import (
    "github.com/spf13/cobra"
    "github.com/pkg/errors"
    "github.com/rzajac/iotdet/pkg/hq"
)

var configureCmd = &cobra.Command{
    Use:   "configure",
    Short: "Configure new agent.",
    Long:  `Configure new agent.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        if len(args) != 1 {
            return errors.New("missing access point name argument")
        }

        ap := hq.NewAgentAP(args[0])
        itf, err := hq.GetInterface(cfg)
        if err != nil {
            return err
        }

        return itf.Configure(ap)
    },
}

func init() {
    rootCmd.AddCommand(configureCmd)
}
