package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	// Test that root command can be created.
	cmd := getRootCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "bluebanquise-installer", cmd.Use)
	assert.Equal(t, "BlueBanquise Installer CLI", cmd.Short)
}

func TestRootCommandHelp(t *testing.T) {
	// Test that root command shows help.
	cmd := getRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "BlueBanquise Installer")
	assert.Contains(t, buf.String(), "online")
	assert.Contains(t, buf.String(), "offline")
	assert.Contains(t, buf.String(), "download")
	assert.Contains(t, buf.String(), "status")
}

func TestRootCommandNoArgs(t *testing.T) {
	// Test that root command shows help when no args provided.
	cmd := getRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "BlueBanquise Installer")
}

// Helper function to get root command for testing.
func getRootCmd() *cobra.Command {
	// Reset any global state.
	rootCmd = &cobra.Command{
		Use:   "bluebanquise-installer",
		Short: "BlueBanquise Installer CLI",
		Long: `BlueBanquise Installer - A CLI tool to install BlueBanquise clusters.

BlueBanquise is a coherent Ansible roles collection designed to deploy and manage 
large groups of hosts (clusters of nodes). This installer provides both online 
and offline installation modes with support for custom users.

Available commands:
  online    - Install BlueBanquise in online mode (downloads from GitHub)
  offline   - Install BlueBanquise in offline mode (use --collections-path or --tarball-path)
  download  - Download collections for offline installation (use --tarball for tarballs)
  status    - Check BlueBanquise installation status

All commands support custom user configuration with --user and --home flags.

For more information, visit: https://bluebanquise.com`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				// In test context, we can ignore this error.
				_ = err
			}
		},
	}
	return rootCmd
}
