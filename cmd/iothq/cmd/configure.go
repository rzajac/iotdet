package cmd

import (
    "github.com/spf13/cobra"
    "github.com/pkg/errors"
)

var configureCmd = &cobra.Command{
    Use:   "configure",
    Short: "configure new agent.",
    Long:  `configure new agent.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        if len(args) != 1 {
            return errors.New("missing access point name argument")
        }

        h, err := getHQ()
        if err != nil {
            return err
        }

        //ap := hq.NewAgentAP(args[0])
        //itf, err := hq.getInterface(cfg)
        //if err != nil {
        //    return err
        //}
        //
        //return itf.Configure(ap)

        return nil
    },
}

func init() {
    rootCmd.AddCommand(configureCmd)
}
