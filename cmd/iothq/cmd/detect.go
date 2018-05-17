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
        cfg, err := config()
        if err != nil {
            return err
        }

        aps, err := hq.Detect(cfg)
        if err != nil {
            return err
        }

        c, err := hq.NewMQTTClient(cfg)
        if err != nil {
            return err
        }

        for _, ap := range aps {
            agent := fmt.Sprintf("AP: %s BSSID: %s", ap.Name, ap.Bssid)
            fmt.Println(agent + "\n")
            token := c.Publish("test/topic", 0, false, agent)
            token.Wait()
        }

        return nil
    },
}

func init() {
    rootCmd.AddCommand(detectCmd)
}
