package cmd

import (
    "github.com/spf13/cobra"
    "os"
    "github.com/rzajac/iotdet/pkg/hq"
)

var detectCmd = &cobra.Command{
    Use:   "detect",
    Short: "Detect new agents.",
    Long:  `Detect new agents.`,
    Run: func(cmd *cobra.Command, args []string) {
        cfg, err := config()
        if err != nil {
            log.Error(err)
            os.Exit(1)
        }

        detector, err := hq.NewDetector(cfg)
        if err != nil {
            log.Error(err)
            os.Exit(1)
        }

        if _, err := detector.Detect(); err != nil {
            log.Error(err)
            os.Exit(1)
        }
    },
}

func init() {
    rootCmd.AddCommand(detectCmd)
}
