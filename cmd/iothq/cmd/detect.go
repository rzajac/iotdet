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

        detector, err := hq.NewDetector(cfg)
        if err != nil {
            return err
        }

        aps, err := detector.Detect()
        if err != nil {
            return err
        }

        for _, ap := range aps {
            fmt.Printf("AP: %s BSSID: %s\n", ap.Name, ap.Bssid)
        }

        return nil
    },
}

func init() {
    rootCmd.AddCommand(detectCmd)
}
