package cmd

import (
    "github.com/spf13/cobra"
    "fmt"
    "os"
)

var detectCmd = &cobra.Command{
    Use:   "detect",
    Short: "Detect new agents.",
    Long:  `Detect new agents.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        h, err := getHQ()
        if err != nil {
            return err
        }

        aps, err := h.DetectAgents()
        if err != nil {
            return err
        }

        // Nothing found.
        if len(aps) == 0 {
            return nil
        }

        for _, ap := range aps {
            fmt.Printf("found new agent: %s\n", ap.Name)
            if err := h.PublishMQTT("hq/new_agent", 0, false, ap.MAC()); err != nil {
               fmt.Fprint(os.Stderr, err)
            }
        }

        return nil
    },
}

func init() {
    rootCmd.AddCommand(detectCmd)
}
