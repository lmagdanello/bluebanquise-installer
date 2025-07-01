package cmd

import (
	"os"

	"github.com/lmagdanello/bluebanquise-installer/internal/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bluebanquise-installer",
	Short: "BlueBanquise Installer CLI",
	Long: `BlueBanquise Installer - A CLI tool to install BlueBanquise clusters.

BlueBanquise is a coherent Ansible roles collection designed to deploy and manage 
large groups of hosts (clusters of nodes). This installer provides both online 
and offline installation modes with support for custom users.

Available commands:
  online    - Install BlueBanquise in online mode (downloads from GitHub)
  offline   - Install BlueBanquise in offline mode (use --collections-path)
  download  - Download collections for offline installation
  status    - Check BlueBanquise installation status

All commands support custom user configuration with --user and --home flags.

For more information, visit: https://bluebanquise.com`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.LogInfo("Showing help information")
		if err := cmd.Help(); err != nil {
			utils.LogError("Error showing help", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.LogError("Root command execution failed", err)
		os.Exit(1)
	}
}
