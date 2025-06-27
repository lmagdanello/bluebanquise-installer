package bootstrap

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lmagdanello/bluebanquise-installer/internal/system"
	"github.com/lmagdanello/bluebanquise-installer/internal/utils"
)

// ConfigureEnvironment sets up the BlueBanquise Python virtual environment and required env vars.
func ConfigureEnvironment(userName, userHome, collectionsPath string) error {
	utils.LogInfo("Configuring BlueBanquise environment", "user", userName, "home", userHome)

	venvDir := filepath.Join(userHome, "ansible_venv")
	bashrc := filepath.Join(userHome, ".bashrc")

	osID, version, err := system.DetectOS()
	if err != nil {
		utils.LogError("Failed to detect OS", err)
		return fmt.Errorf("failed to detect OS: %v", err)
	}
	utils.LogInfo("OS detected for environment configuration", "os", osID, "version", version)

	// RHEL7 specific: Export rh-python38
	if osID == "rhel" && version == "7" {
		utils.LogInfo("Configuring RHEL7 specific environment")
		if err := utils.ExportRHPython38(userHome); err != nil {
			utils.LogError("Failed to export rh-python38 environment", err)
			return fmt.Errorf("failed to export rh-python38 environment: %v", err)
		}
	}

	utils.LogInfo("Creating Python virtual environment", "path", venvDir)
	fmt.Println("Creating Python virtual environment...")
	if err := utils.InstallPackages([]string{"python3", "python3-venv"}); err != nil {
		utils.LogError("Failed to install python3", err)
		return fmt.Errorf("failed to install python3: %v", err)
	}
	if _, err := os.Stat(venvDir); os.IsNotExist(err) {
		if err := utils.InstallPackages([]string{"python3-venv"}); err != nil {
			utils.LogError("Failed to install python3-venv", err)
			return fmt.Errorf("failed to venv package missing: %v", err)
		}
	}

	cmd := filepath.Join("/usr/bin", "python3")
	utils.LogCommand(cmd, "-m", "venv", venvDir)
	if err := utils.RunCommand(cmd, "-m", "venv", venvDir); err != nil {
		utils.LogError("Failed to create virtualenv", err, "path", venvDir)
		return fmt.Errorf("failed to create virtualenv: %v", err)
	}

	utils.LogInfo("Installing Python requirements", "requirements", system.PythonRequirements)
	if err := utils.InstallRequirements(venvDir, system.PythonRequirements); err != nil {
		utils.LogError("Failed to install Python packages", err, "venv", venvDir)
		return fmt.Errorf("failed to install Python packages: %v", err)
	}

	// Add to .bashrc
	utils.LogInfo("Updating .bashrc with environment variables", "file", bashrc)
	exportLines := []string{
		fmt.Sprintf("source %s", filepath.Join(venvDir, "bin", "activate")),
		"export ANSIBLE_CONFIG=$HOME/bluebanquise/ansible.cfg",
	}
	for _, line := range exportLines {
		if err := utils.AppendLineIfMissing(bashrc, line); err != nil {
			utils.LogError("Failed to update .bashrc", err, "line", line)
			return fmt.Errorf("failed to update .bashrc: %v", err)
		}
	}

	// Ensure sudoers has PYTHONPATH preserved
	utils.LogInfo("Updating sudoers to preserve PYTHONPATH")
	if err := utils.EnsureLineInSudoers(`Defaults env_keep += "PYTHONPATH"`); err != nil {
		utils.LogError("Failed to update sudoers", err)
		return fmt.Errorf("failed to update sudoers: %v", err)
	}

	// Configure SSH
	utils.LogInfo("Configuring SSH", "home", userHome)
	fmt.Println("Configuring SSH...")
	if err := utils.ConfigureSSH(userHome); err != nil {
		utils.LogError("Failed to configure SSH", err, "home", userHome)
		return fmt.Errorf("failed to configure SSH: %v", err)
	}

	// Create bluebanquise directory for ansible.cfg
	bluebanquiseDir := filepath.Join(userHome, "bluebanquise")
	utils.LogInfo("Creating bluebanquise directory", "path", bluebanquiseDir)
	if err := os.MkdirAll(bluebanquiseDir, 0755); err != nil {
		utils.LogError("Failed to create bluebanquise directory", err, "path", bluebanquiseDir)
		return fmt.Errorf("failed to create bluebanquise directory: %v", err)
	}

	utils.LogInfo("Environment configured successfully", "user", userName, "home", userHome)
	fmt.Println("Environment configured successfully.")
	return nil
}

// ConfigureEnvironmentOffline sets up the BlueBanquise Python virtual environment using offline requirements.
func ConfigureEnvironmentOffline(userName, userHome, requirementsPath string) error {
	utils.LogInfo("Configuring BlueBanquise environment offline", "user", userName, "home", userHome, "requirements_path", requirementsPath)

	// Detect OS and configure RHEL7 specific settings
	if err := configureOSSpecificSettings(userHome); err != nil {
		return err
	}

	// Create virtual environment
	venvDir := filepath.Join(userHome, "ansible_venv")
	if err := createVirtualEnvironment(venvDir); err != nil {
		return err
	}

	// Install requirements offline if path provided
	if err := installOfflineRequirements(venvDir, requirementsPath); err != nil {
		return err
	}

	// Configure environment files
	if err := configureEnvironmentFiles(userHome, venvDir); err != nil {
		return err
	}

	utils.LogInfo("Offline environment configured successfully", "user", userName, "home", userHome, "requirements_path", requirementsPath)
	fmt.Println("Environment configured successfully.")
	return nil
}

// configureOSSpecificSettings handles OS-specific configuration like RHEL7 rh-python38.
func configureOSSpecificSettings(userHome string) error {
	osID, version, err := system.DetectOS()
	if err != nil {
		utils.LogError("Failed to detect OS", err)
		return fmt.Errorf("failed to detect OS: %v", err)
	}
	utils.LogInfo("OS detected for offline environment configuration", "os", osID, "version", version)

	// RHEL7 specific: Export rh-python38
	if osID == "rhel" && version == "7" {
		utils.LogInfo("Configuring RHEL7 specific environment")
		if err := utils.ExportRHPython38(userHome); err != nil {
			utils.LogError("Failed to export rh-python38 environment", err)
			return fmt.Errorf("failed to export rh-python38 environment: %v", err)
		}
	}

	return nil
}

// createVirtualEnvironment creates the Python virtual environment.
func createVirtualEnvironment(venvDir string) error {
	utils.LogInfo("Creating Python virtual environment", "path", venvDir)
	fmt.Println("Creating Python virtual environment...")

	if err := utils.InstallPackages([]string{"python3", "python3-venv"}); err != nil {
		utils.LogError("Failed to install python3", err)
		return fmt.Errorf("failed to install python3: %v", err)
	}

	if _, err := os.Stat(venvDir); os.IsNotExist(err) {
		if err := utils.InstallPackages([]string{"python3-venv"}); err != nil {
			utils.LogError("Failed to install python3-venv", err)
			return fmt.Errorf("failed to venv package missing: %v", err)
		}
	}

	cmd := filepath.Join("/usr/bin", "python3")
	utils.LogCommand(cmd, "-m", "venv", venvDir)
	if err := utils.RunCommand(cmd, "-m", "venv", venvDir); err != nil {
		utils.LogError("Failed to create virtualenv", err, "path", venvDir)
		return fmt.Errorf("failed to create virtualenv: %v", err)
	}

	return nil
}

// installOfflineRequirements installs Python requirements from offline path.
func installOfflineRequirements(venvDir, requirementsPath string) error {
	if requirementsPath != "" {
		utils.LogInfo("Installing Python requirements offline", "requirements_path", requirementsPath)
		if err := utils.InstallRequirementsOffline(venvDir, requirementsPath); err != nil {
			utils.LogError("Failed to install Python packages offline", err, "venv", venvDir, "requirements_path", requirementsPath)
			return fmt.Errorf("failed to install Python packages offline: %v", err)
		}
	} else {
		utils.LogInfo("No requirements path provided, skipping Python package installation")
	}
	return nil
}

// configureEnvironmentFiles sets up .bashrc, sudoers, SSH, and bluebanquise directory.
func configureEnvironmentFiles(userHome, venvDir string) error {
	bashrc := filepath.Join(userHome, ".bashrc")

	// Add to .bashrc
	utils.LogInfo("Updating .bashrc with environment variables", "file", bashrc)
	exportLines := []string{
		fmt.Sprintf("source %s", filepath.Join(venvDir, "bin", "activate")),
		"export ANSIBLE_CONFIG=$HOME/bluebanquise/ansible.cfg",
	}
	for _, line := range exportLines {
		if err := utils.AppendLineIfMissing(bashrc, line); err != nil {
			utils.LogError("Failed to update .bashrc", err, "line", line)
			return fmt.Errorf("failed to update .bashrc: %v", err)
		}
	}

	// Ensure sudoers has PYTHONPATH preserved
	utils.LogInfo("Updating sudoers to preserve PYTHONPATH")
	if err := utils.EnsureLineInSudoers(`Defaults env_keep += "PYTHONPATH"`); err != nil {
		utils.LogError("Failed to update sudoers", err)
		return fmt.Errorf("failed to update sudoers: %v", err)
	}

	// Configure SSH
	utils.LogInfo("Configuring SSH", "home", userHome)
	fmt.Println("Configuring SSH...")
	if err := utils.ConfigureSSH(userHome); err != nil {
		utils.LogError("Failed to configure SSH", err, "home", userHome)
		return fmt.Errorf("failed to configure SSH: %v", err)
	}

	// Create bluebanquise directory for ansible.cfg
	bluebanquiseDir := filepath.Join(userHome, "bluebanquise")
	utils.LogInfo("Creating bluebanquise directory", "path", bluebanquiseDir)
	if err := os.MkdirAll(bluebanquiseDir, 0755); err != nil {
		utils.LogError("Failed to create bluebanquise directory", err, "path", bluebanquiseDir)
		return fmt.Errorf("failed to create bluebanquise directory: %v", err)
	}

	return nil
}
