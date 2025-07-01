package bootstrap

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/lmagdanello/bluebanquise-installer/internal/system"
	"github.com/lmagdanello/bluebanquise-installer/internal/utils"
)

// InstallCollectionsOnline installs BlueBanquise collections from GitHub.
func InstallCollectionsOnline(userHome string) error {
	utils.LogInfo("Installing collections online", "home", userHome)

	venvDir := filepath.Join(userHome, "ansible_venv")
	venvBin := filepath.Join(venvDir, "bin")
	ansibleGalaxy := filepath.Join(venvBin, "ansible-galaxy")
	collectionsDir := filepath.Join(userHome, ".ansible", "collections")

	// Verify ansible-galaxy exists, create environment if it doesn't
	if err := ensureAnsibleGalaxy(venvDir, ansibleGalaxy); err != nil {
		return err
	}

	// Create collections directory if it doesn't exist.
	if err := os.MkdirAll(collectionsDir, 0755); err != nil {
		utils.LogError("Failed to create collections directory", err, "path", collectionsDir)
		return fmt.Errorf("failed to create collections directory: %v", err)
	}

	utils.LogInfo("Installing BlueBanquise collections", "collections_dir", collectionsDir)
	fmt.Println("Installing BlueBanquise collections...")

	utils.LogCommand(ansibleGalaxy, "collection", "install", "git+https://github.com/bluebanquise/bluebanquise.git#/collections/infrastructure,master", "-p", collectionsDir)
	cmd := exec.Command(ansibleGalaxy, "collection", "install", "git+https://github.com/bluebanquise/bluebanquise.git#/collections/infrastructure,master", "-p", collectionsDir)
	if err := cmd.Run(); err != nil {
		utils.LogError("Failed to install BlueBanquise collections", err)
		return fmt.Errorf("failed to install BlueBanquise collections: %v", err)
	}

	utils.LogInfo("Installing community.general collection", "collections_dir", collectionsDir)
	fmt.Println("Installing community.general collection...")

	utils.LogCommand(ansibleGalaxy, "collection", "install", "community.general", "-p", collectionsDir)
	cmd = exec.Command(ansibleGalaxy, "collection", "install", "community.general", "-p", collectionsDir)
	if err := cmd.Run(); err != nil {
		utils.LogError("Failed to install community.general collection", err)
		return fmt.Errorf("failed to install community.general collection: %v", err)
	}

	utils.LogInfo("Collections installed successfully online", "collections_dir", collectionsDir)
	return nil
}

// InstallCollectionsFromPath installs BlueBanquise collections from a given path.
func InstallCollectionsFromPath(path, userHome string) error {
	utils.LogInfo("Installing collections from path", "path", path, "home", userHome)
	venvDir := filepath.Join(userHome, "ansible_venv")
	venvBin := filepath.Join(venvDir, "bin")
	ansibleGalaxy := filepath.Join(venvBin, "ansible-galaxy")
	collectionsDir := filepath.Join(userHome, ".ansible", "collections")

	// Verify ansible-galaxy exists, create environment if it doesn't
	if err := ensureAnsibleGalaxy(venvDir, ansibleGalaxy); err != nil {
		return err
	}

	// Create collections directory if it doesn't exist.
	if err := os.MkdirAll(collectionsDir, 0755); err != nil {
		utils.LogError("Failed to create collections directory", err, "path", collectionsDir)
		return fmt.Errorf("failed to create collections directory: %v", err)
	}
	// Check if path is a file or directory.
	info, err := os.Stat(path)
	if err != nil {
		utils.LogError("Failed to stat path", err, "path", path)
		return fmt.Errorf("failed to stat path: %v", err)
	}
	if info.IsDir() {
		// Directory containing multiple tarballs/collections.
		utils.LogInfo("Processing directory", "path", path)
		entries, err := os.ReadDir(path)
		if err != nil {
			utils.LogError("Failed to read directory", err, "path", path)
			return fmt.Errorf("failed to read directory: %v", err)
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				name := entry.Name()
				if strings.HasSuffix(name, ".tar.gz") || strings.HasSuffix(name, ".tgz") {
					file := filepath.Join(path, name)
					utils.LogInfo("Installing collection from file", "file", name, "path", file)
					fmt.Printf("Installing collection from file: %s\n", name)
					utils.LogCommand(ansibleGalaxy, "collection", "install", file, "-p", collectionsDir)
					cmd := exec.Command(ansibleGalaxy, "collection", "install", file, "-p", collectionsDir)
					if err := cmd.Run(); err != nil {
						utils.LogError("Failed to install collection from file", err, "file", name, "path", file)
						return fmt.Errorf("failed to install collection from file %s: %v", name, err)
					}
				}
			}
		}
	} else {
		// Single file.
		utils.LogInfo("Installing collection from single file", "file", filepath.Base(path), "path", path)
		fmt.Printf("Installing collection from file: %s\n", filepath.Base(path))
		utils.LogCommand(ansibleGalaxy, "collection", "install", path, "-p", collectionsDir)
		cmd := exec.Command(ansibleGalaxy, "collection", "install", path, "-p", collectionsDir)
		if err := cmd.Run(); err != nil {
			utils.LogError("Failed to install collection from file", err, "path", path)
			return fmt.Errorf("failed to install collection from file: %v", err)
		}
	}
	utils.LogInfo("Collections installed successfully from path", "path", path)
	return nil
}

// InstallCoreVariablesOnline installs core variables by downloading from GitHub.
func InstallCoreVariablesOnline(userHome string) error {
	utils.LogInfo("Installing core variables online", "home", userHome)

	// Validate userHome is not empty.
	if userHome == "" {
		utils.LogError("User home directory is empty", nil)
		return fmt.Errorf("user home directory cannot be empty")
	}

	// Create inventory directory structure.
	inventoryDir := filepath.Join(userHome, "bluebanquise", "inventory")
	groupVarsDir := filepath.Join(inventoryDir, "group_vars", "all")

	utils.LogInfo("Creating inventory directory structure", "path", groupVarsDir)
	if err := os.MkdirAll(groupVarsDir, 0755); err != nil {
		utils.LogError("Failed to create inventory directory", err, "path", groupVarsDir)
		return fmt.Errorf("failed to create inventory directory: %v", err)
	}

	// Download bb_core.yml from GitHub.
	bbCoreURL := "https://raw.githubusercontent.com/bluebanquise/bluebanquise/refs/heads/master/resources/bb_core.yml"
	bbCorePath := filepath.Join(groupVarsDir, "bb_core.yml")

	utils.LogInfo("Downloading bb_core.yml", "url", bbCoreURL, "path", bbCorePath)
	fmt.Println("Downloading core variables from GitHub...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", bbCoreURL, http.NoBody)
	if err != nil {
		utils.LogError("Failed to create request", err, "url", bbCoreURL)
		return fmt.Errorf("failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.LogError("Failed to download bb_core.yml", err, "url", bbCoreURL)
		return fmt.Errorf("failed to download bb_core.yml: %v", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			utils.LogWarning("Failed to close response body", "error", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		utils.LogError("Failed to download bb_core.yml", nil, "status", resp.StatusCode, "url", bbCoreURL)
		return fmt.Errorf("failed to download bb_core.yml: HTTP %d", resp.StatusCode)
	}

	file, err := os.Create(bbCorePath)
	if err != nil {
		utils.LogError("Failed to create bb_core.yml file", err, "path", bbCorePath)
		return fmt.Errorf("failed to create bb_core.yml file: %v", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			utils.LogWarning("Failed to close file", "error", closeErr)
		}
	}()

	if _, err := io.Copy(file, resp.Body); err != nil {
		utils.LogError("Failed to write bb_core.yml file", err, "path", bbCorePath)
		return fmt.Errorf("failed to write bb_core.yml file: %v", err)
	}

	utils.LogInfo("Core variables installed successfully online", "path", bbCorePath)
	fmt.Println("Core variables installed successfully.")
	return nil
}

// InstallCoreVariablesOffline installs core variables from local path.
func InstallCoreVariablesOffline(coreVarsPath, userHome string) error {
	utils.LogInfo("Installing core variables offline", "core_vars_path", coreVarsPath, "home", userHome)

	// Validate userHome is not empty.
	if userHome == "" {
		utils.LogError("User home directory is empty", nil)
		return fmt.Errorf("user home directory cannot be empty")
	}

	// Create inventory directory structure.
	inventoryDir := filepath.Join(userHome, "bluebanquise", "inventory")
	groupVarsDir := filepath.Join(inventoryDir, "group_vars", "all")

	utils.LogInfo("Creating inventory directory structure", "path", groupVarsDir)
	if err := os.MkdirAll(groupVarsDir, 0755); err != nil {
		utils.LogError("Failed to create inventory directory", err, "path", groupVarsDir)
		return fmt.Errorf("failed to create inventory directory: %v", err)
	}

	// Check if coreVarsPath is a file or directory.
	info, err := os.Stat(coreVarsPath)
	if err != nil {
		utils.LogError("Failed to stat core variables path", err, "path", coreVarsPath)
		return fmt.Errorf("failed to stat core variables path: %v", err)
	}

	if info.IsDir() {
		// Directory containing multiple variable files.
		utils.LogInfo("Processing core variables directory", "path", coreVarsPath)
		entries, err := os.ReadDir(coreVarsPath)
		if err != nil {
			utils.LogError("Failed to read core variables directory", err, "path", coreVarsPath)
			return fmt.Errorf("failed to read core variables directory: %v", err)
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				name := entry.Name()
				if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
					sourceFile := filepath.Join(coreVarsPath, name)
					destFile := filepath.Join(groupVarsDir, name)

					utils.LogInfo("Installing core variable file", "file", name, "source", sourceFile, "dest", destFile)
					fmt.Printf("Installing core variable file: %s\n", name)

					if err := copyFile(sourceFile, destFile); err != nil {
						utils.LogError("Failed to copy core variable file", err, "file", name, "source", sourceFile)
						return fmt.Errorf("failed to copy core variable file %s: %v", name, err)
					}
				}
			}
		}
	} else {
		// Single variable file.
		destFile := filepath.Join(groupVarsDir, "bb_core.yml")
		utils.LogInfo("Installing core variable file", "source", coreVarsPath, "dest", destFile)
		fmt.Printf("Installing core variable file: %s\n", filepath.Base(coreVarsPath))

		if err := copyFile(coreVarsPath, destFile); err != nil {
			utils.LogError("Failed to copy core variable file", err, "source", coreVarsPath, "dest", destFile)
			return fmt.Errorf("failed to copy core variable file: %v", err)
		}
	}

	utils.LogInfo("Core variables installed successfully offline", "path", coreVarsPath)
	fmt.Println("Core variables installed successfully.")
	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := sourceFile.Close(); closeErr != nil {
			utils.LogWarning("Failed to close source file", "error", closeErr)
		}
	}()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := destFile.Close(); closeErr != nil {
			utils.LogWarning("Failed to close destination file", "error", closeErr)
		}
	}()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// ensureAnsibleGalaxy ensures that ansible-galaxy is available in the virtual environment.
func ensureAnsibleGalaxy(venvDir, ansibleGalaxy string) error {
	if _, err := os.Stat(ansibleGalaxy); os.IsNotExist(err) {
		utils.LogInfo("ansible-galaxy not found, creating environment", "path", ansibleGalaxy)
		fmt.Println("Creating Python environment for collections installation...")

		// Create virtual environment
		if err := createVirtualEnvironment(venvDir); err != nil {
			return fmt.Errorf("failed to create virtual environment: %v", err)
		}

		// Install requirements to get ansible-galaxy
		utils.LogInfo("Installing Python requirements for ansible-galaxy", "requirements", system.PythonRequirements)
		if err := utils.InstallRequirements(venvDir, system.PythonRequirements); err != nil {
			utils.LogError("Failed to install Python packages", err, "venv", venvDir)
			return fmt.Errorf("failed to install Python packages: %v", err)
		}

		// Verify ansible-galaxy exists now
		if _, err := os.Stat(ansibleGalaxy); os.IsNotExist(err) {
			utils.LogError("ansible-galaxy still not found after environment setup", err, "path", ansibleGalaxy)
			return fmt.Errorf("ansible-galaxy not found at %s after environment setup: %v", ansibleGalaxy, err)
		}
	}
	return nil
}
