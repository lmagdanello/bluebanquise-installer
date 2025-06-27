package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lmagdanello/bluebanquise-installer/internal/bootstrap"
	"github.com/lmagdanello/bluebanquise-installer/internal/utils"
	"github.com/spf13/cobra"
)

var (
	statusUserName string
	statusCmd      = &cobra.Command{
		Use:   "status",
		Short: "Check BlueBanquise installation status",
		Long: `Check the status of BlueBanquise installation.

This command verifies:
- User existence and home directory
- Python virtual environment
- Ansible installation
- BlueBanquise collections
- Core variables

Examples:
  # Check status for default user (bluebanquise)
  ./bluebanquise-installer status

  # Check status for specific user
  ./bluebanquise-installer status --user myuser`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := checkStatus(); err != nil {
				utils.LogError("Status check failed", err)
				fmt.Printf("Status check failed: %v\n", err)
				os.Exit(1)
			}
		},
	}
)

func checkStatus() error {
	utils.LogInfo("Checking BlueBanquise installation status", "user", statusUserName)

	// Check user and home directory
	userHome, err := getUserHome(statusUserName)
	if err != nil {
		return fmt.Errorf("%s user home directory not found", statusUserName)
	}

	fmt.Printf("✓ User %s home directory: %s\n", statusUserName, userHome)

	// Check Python virtual environment
	venvDir := filepath.Join(userHome, "ansible_venv")
	if _, err := os.Stat(venvDir); os.IsNotExist(err) {
		return fmt.Errorf("python virtual environment not found")
	}

	fmt.Printf("✓ Python virtual environment: %s\n", venvDir)

	// Check if activate script exists
	activateScript := filepath.Join(venvDir, "bin", "activate")
	if _, err := os.Stat(activateScript); os.IsNotExist(err) {
		return fmt.Errorf("virtual environment activate script not found")
	}

	// Check Ansible installation
	ansiblePath := filepath.Join(venvDir, "bin", "ansible")
	if _, err := os.Stat(ansiblePath); os.IsNotExist(err) {
		return fmt.Errorf("ansible not found in virtual environment")
	}

	fmt.Printf("✓ Ansible: %s\n", ansiblePath)

	ansibleGalaxyPath := filepath.Join(venvDir, "bin", "ansible-galaxy")
	if _, err := os.Stat(ansibleGalaxyPath); os.IsNotExist(err) {
		return fmt.Errorf("ansible-galaxy not found in virtual environment")
	}

	fmt.Printf("✓ Ansible Galaxy: %s\n", ansibleGalaxyPath)

	// Check BlueBanquise collections
	collectionsDir := filepath.Join(userHome, ".ansible", "collections")
	if _, err := os.Stat(collectionsDir); os.IsNotExist(err) {
		return fmt.Errorf("bluebanquise collections not found")
	}

	fmt.Printf("✓ Collections directory: %s\n", collectionsDir)

	// Check if infrastructure collection exists
	infraCollectionDir := filepath.Join(collectionsDir, "ansible_collections", "bluebanquise", "infrastructure")
	if _, err := os.Stat(infraCollectionDir); os.IsNotExist(err) {
		return fmt.Errorf("bluebanquise infrastructure collection not found")
	}

	fmt.Printf("✓ BlueBanquise infrastructure collection: %s\n", infraCollectionDir)

	// Check core variables
	coreVarsPath := filepath.Join(userHome, "bluebanquise", "inventory", "group_vars", "all", "bb_core.yml")
	if _, err := os.Stat(coreVarsPath); os.IsNotExist(err) {
		fmt.Printf("⚠ Core variables not found: %s\n", coreVarsPath)
	} else {
		fmt.Printf("✓ Core variables: %s\n", coreVarsPath)
	}

	utils.LogInfo("BlueBanquise installation status check completed successfully", "user", statusUserName)
	fmt.Println("\n✓ BlueBanquise installation is ready!")
	return nil
}

func getUserHome(userName string) (string, error) {
	if userName == "" {
		userName = "bluebanquise"
	}

	_, _, err := bootstrap.GetUserInfo(userName)
	if err != nil {
		return "", err
	}

	// Get home directory from /etc/passwd or use default
	homeDir := fmt.Sprintf("/home/%s", userName)
	if _, err := os.Stat(homeDir); os.IsNotExist(err) {
		// Try alternative locations
		altDirs := []string{
			fmt.Sprintf("/var/lib/%s", userName),
			fmt.Sprintf("/opt/%s", userName),
		}
		for _, dir := range altDirs {
			if _, err := os.Stat(dir); err == nil {
				homeDir = dir
				break
			}
		}
	}

	return homeDir, nil
}

func init() {
	statusCmd.Flags().StringVarP(&statusUserName, "user", "u", "", "Username to check status for (default: bluebanquise)")
	rootCmd.AddCommand(statusCmd)
}
