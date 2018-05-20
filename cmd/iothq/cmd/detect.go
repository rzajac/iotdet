package cmd

import (
    "github.com/spf13/cobra"
    "fmt"
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

        agents, err := h.DetectAgents()
        if err != nil {
            return err
        }

        for _, agent := range agents {
            fmt.Printf("found new agent: %s (%s)\n", agent.Name, agent.MAC())
        }

        return nil
    },
}

func init() {
    rootCmd.AddCommand(detectCmd)
}
