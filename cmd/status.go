package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lmagdanello/bluebanquise-installer/internal/utils"
	"github.com/spf13/cobra"
)

var (
	statusUserName string
	statusUserHome string
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check BlueBanquise installation status",
	Long: `Check the status of BlueBanquise installation.
	
This command will verify:
1. BlueBanquise user existence
2. Python virtual environment
3. Ansible installation
4. Collections installation`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.LogInfo("Checking BlueBanquise installation status",
			"user", statusUserName,
			"home", statusUserHome)

		checks := []struct {
			name  string
			check func() error
		}{
			{"BlueBanquise user", checkBlueBanquiseUser},
			{"Python virtual environment", checkPythonVenv},
			{"Ansible installation", checkAnsibleInstallation},
			{"BlueBanquise collections", checkBlueBanquiseCollections},
		}

		allOK := true
		for _, c := range checks {
			utils.LogInfo("Running check", "check", c.name)
			fmt.Printf("Checking %s... ", c.name)
			if err := c.check(); err != nil {
				utils.LogError("Check failed", err, "check", c.name)
				fmt.Printf("FAILED: %v\n", err)
				allOK = false
			} else {
				utils.LogInfo("Check passed", "check", c.name)
				fmt.Println("OK")
			}
		}

		if allOK {
			utils.LogInfo("All checks passed - BlueBanquise is properly installed")
			fmt.Println("\nBlueBanquise is properly installed and configured!")
		} else {
			utils.LogError("Some checks failed - BlueBanquise installation has issues", nil)
			fmt.Println("\nBlueBanquise installation has issues. Please reinstall.")
			os.Exit(1)
		}
	},
}

func checkBlueBanquiseUser() error {
	if _, err := os.Stat(statusUserHome); os.IsNotExist(err) {
		return fmt.Errorf("%s user home directory not found", statusUserName)
	}
	return nil
}

func checkPythonVenv() error {
	venvPath := filepath.Join(statusUserHome, "ansible_venv")
	if _, err := os.Stat(venvPath); os.IsNotExist(err) {
		return fmt.Errorf("Python virtual environment not found")
	}

	// Check if activate script exists
	activateScript := filepath.Join(venvPath, "bin", "activate")
	if _, err := os.Stat(activateScript); os.IsNotExist(err) {
		return fmt.Errorf("virtual environment activate script not found")
	}

	return nil
}

func checkAnsibleInstallation() error {
	ansiblePath := filepath.Join(statusUserHome, "ansible_venv", "bin", "ansible")
	if _, err := os.Stat(ansiblePath); os.IsNotExist(err) {
		return fmt.Errorf("ansible not found in virtual environment")
	}

	ansibleGalaxyPath := filepath.Join(statusUserHome, "ansible_venv", "bin", "ansible-galaxy")
	if _, err := os.Stat(ansibleGalaxyPath); os.IsNotExist(err) {
		return fmt.Errorf("ansible-galaxy not found in virtual environment")
	}

	return nil
}

func checkBlueBanquiseCollections() error {
	collectionsPath := filepath.Join(statusUserHome, ".ansible", "collections", "ansible_collections", "bluebanquise")
	if _, err := os.Stat(collectionsPath); os.IsNotExist(err) {
		return fmt.Errorf("BlueBanquise collections not found")
	}

	// Check if infrastructure collection exists
	infraPath := filepath.Join(collectionsPath, "infrastructure")
	if _, err := os.Stat(infraPath); os.IsNotExist(err) {
		return fmt.Errorf("BlueBanquise infrastructure collection not found")
	}

	return nil
}

func init() {
	statusCmd.Flags().StringVarP(&statusUserName, "user", "u", "bluebanquise", "Username for BlueBanquise")
	statusCmd.Flags().StringVarP(&statusUserHome, "home", "h", "/var/lib/bluebanquise", "Home directory for BlueBanquise user")

	rootCmd.AddCommand(statusCmd)
}
