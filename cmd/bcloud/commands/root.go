package commands

import (
	"fmt"

	"github.com/brevdev/cloud/cmd/bcloud/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "bcloud <refId> <verb>",
	Short: "Brev Cloud CLI for managing GPU compute across providers",
	Long: `A vendor-agnostic CLI for managing clusterable, GPU-accelerated compute
across multiple cloud providers using the Brev Cloud SDK.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.bcloud/credentials.yaml)")

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(terminateCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(typesCmd)
}

func initConfig() {
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		fmt.Println("Create ~/.bcloud/credentials.yaml with your cloud credentials")
	}
}
