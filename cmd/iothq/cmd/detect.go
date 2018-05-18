package cmd

import (
    "github.com/spf13/cobra"
    "github.com/rzajac/iotdet/pkg/hq"
    "fmt"
)

var detectCmd = &cobra.Command{
    Use:   "detect",
    Short: "Detect new agents.",
    Long:  `Detect new agents.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        aps, err := hq.Detect(cfg)
        if err != nil {
            return err
        }

        // Nothing found.
        if len(aps) == 0 {
            return nil
        }

        c, err := hq.NewMQTTClient(cfg)
        if err != nil {
            return err
        }

        for _, ap := range aps {
            fmt.Printf("found new agent: %s\n", ap.Name)
            c.Publish("hq/new_agent", 0, false, ap.MAC()).Wait()
        }

        return nil
    },
}

func init() {
    rootCmd.AddCommand(detectCmd)
}
