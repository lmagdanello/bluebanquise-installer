package bootstrap

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lmagdanello/bluebanquise-installer/internal/utils"
)

func InstallCollectionsOnline(userHome string) error {
	utils.LogInfo("Installing collections online", "home", userHome)

	venvBin := filepath.Join(userHome, "ansible_venv", "bin")
	ansibleGalaxy := filepath.Join(venvBin, "ansible-galaxy")

	utils.LogInfo("Installing community.general collection")
	fmt.Println("Installing community.general collection...")
	for i := range 3 {
		utils.LogCommand(ansibleGalaxy, "collection", "install", "community.general")
		cmd := exec.Command(ansibleGalaxy, "collection", "install", "community.general")
		err := cmd.Run()
		if err == nil {
			utils.LogInfo("community.general collection installed successfully")
			break
		}
		utils.LogWarning("community.general installation failed, retrying", "attempt", i+1, "error", err)
		fmt.Println("Retrying community.general installation...")
	}

	utils.LogInfo("Installing BlueBanquise collection from GitHub")
	fmt.Println("Installing BlueBanquise collection from GitHub...")
	utils.LogCommand(ansibleGalaxy, "collection", "install", "git+https://github.com/bluebanquise/bluebanquise.git#/collections/infrastructure,master", "-vvv", "--upgrade")
	cmd := exec.Command(ansibleGalaxy,
		"collection", "install",
		"git+https://github.com/bluebanquise/bluebanquise.git#/collections/infrastructure,master",
		"-vvv", "--upgrade")
	if err := cmd.Run(); err != nil {
		utils.LogError("Failed to install BlueBanquise collection from GitHub", err)
		return fmt.Errorf("failed to install BlueBanquise collection from GitHub: %v", err)
	}

	utils.LogInfo("Collections installed successfully online")
	return nil
}

func InstallCollectionsOffline(collectionsPath, userHome string) error {
	utils.LogInfo("Installing collections offline", "collections_path", collectionsPath, "home", userHome)

	venvBin := filepath.Join(userHome, "ansible_venv", "bin")
	ansibleGalaxy := filepath.Join(venvBin, "ansible-galaxy")
	fmt.Printf("Installing BlueBanquise collection from local path: %s\n", collectionsPath)

	utils.LogCommand(ansibleGalaxy, "collection", "install", collectionsPath)
	cmd := exec.Command(ansibleGalaxy, "collection", "install", collectionsPath)
	if err := cmd.Run(); err != nil {
		utils.LogError("Failed to install local BlueBanquise collection", err, "path", collectionsPath)
		return fmt.Errorf("failed to install local BlueBanquise collection: %v", err)
	}

	utils.LogInfo("Collections installed successfully offline", "path", collectionsPath)
	return nil
}

func InstallCollectionsFromTarballs(tarballPath, userHome string) error {
	utils.LogInfo("Installing collections from tarballs", "tarball_path", tarballPath, "home", userHome)

	venvBin := filepath.Join(userHome, "ansible_venv", "bin")
	ansibleGalaxy := filepath.Join(venvBin, "ansible-galaxy")
	collectionsDir := filepath.Join(userHome, ".ansible", "collections")

	// Create collections directory if it doesn't exist
	if err := os.MkdirAll(collectionsDir, 0755); err != nil {
		utils.LogError("Failed to create collections directory", err, "path", collectionsDir)
		return fmt.Errorf("failed to create collections directory: %v", err)
	}

	// Check if tarballPath is a file or directory
	info, err := os.Stat(tarballPath)
	if err != nil {
		utils.LogError("Failed to stat tarball path", err, "path", tarballPath)
		return fmt.Errorf("failed to stat tarball path: %v", err)
	}

	if info.IsDir() {
		// Directory containing multiple tarballs
		utils.LogInfo("Processing tarball directory", "path", tarballPath)
		entries, err := os.ReadDir(tarballPath)
		if err != nil {
			utils.LogError("Failed to read tarball directory", err, "path", tarballPath)
			return fmt.Errorf("failed to read tarball directory: %v", err)
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				name := entry.Name()
				if strings.HasSuffix(name, ".tar.gz") || strings.HasSuffix(name, ".tgz") {
					tarballFile := filepath.Join(tarballPath, name)
					utils.LogInfo("Installing collection from tarball", "file", name, "path", tarballFile)
					fmt.Printf("Installing collection from tarball: %s\n", name)

					utils.LogCommand(ansibleGalaxy, "collection", "install", tarballFile, "-p", collectionsDir)
					cmd := exec.Command(ansibleGalaxy, "collection", "install", tarballFile, "-p", collectionsDir)
					if err := cmd.Run(); err != nil {
						utils.LogError("Failed to install collection from tarball", err, "file", name, "path", tarballFile)
						return fmt.Errorf("failed to install collection from tarball %s: %v", name, err)
					}
				}
			}
		}
	} else {
		// Single tarball file
		utils.LogInfo("Installing collection from single tarball", "file", filepath.Base(tarballPath), "path", tarballPath)
		fmt.Printf("Installing collection from tarball: %s\n", filepath.Base(tarballPath))

		utils.LogCommand(ansibleGalaxy, "collection", "install", tarballPath, "-p", collectionsDir)
		cmd := exec.Command(ansibleGalaxy, "collection", "install", tarballPath, "-p", collectionsDir)
		if err := cmd.Run(); err != nil {
			utils.LogError("Failed to install collection from tarball", err, "path", tarballPath)
			return fmt.Errorf("failed to install collection from tarball: %v", err)
		}
	}

	utils.LogInfo("Collections installed successfully from tarballs", "path", tarballPath)
	return nil
}

func InstallCoreVariablesOnline(userHome string) error {
	utils.LogInfo("Installing core variables online", "home", userHome)

	// Validate userHome is not empty
	if userHome == "" {
		utils.LogError("User home directory is empty", nil)
		return fmt.Errorf("user home directory cannot be empty")
	}

	// Create inventory directory structure
	inventoryDir := filepath.Join(userHome, "bluebanquise", "inventory")
	groupVarsDir := filepath.Join(inventoryDir, "group_vars", "all")

	utils.LogInfo("Creating inventory directory structure", "path", groupVarsDir)
	if err := os.MkdirAll(groupVarsDir, 0755); err != nil {
		utils.LogError("Failed to create inventory directory", err, "path", groupVarsDir)
		return fmt.Errorf("failed to create inventory directory: %v", err)
	}

	// Download bb_core.yml from GitHub
	bbCoreURL := "https://raw.githubusercontent.com/bluebanquise/bluebanquise/master/resources/bb_core.yml"
	bbCorePath := filepath.Join(groupVarsDir, "bb_core.yml")

	utils.LogInfo("Downloading bb_core.yml", "url", bbCoreURL, "path", bbCorePath)
	fmt.Println("Downloading core variables from GitHub...")

	resp, err := http.Get(bbCoreURL)
	if err != nil {
		utils.LogError("Failed to download bb_core.yml", err, "url", bbCoreURL)
		return fmt.Errorf("failed to download bb_core.yml: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		utils.LogError("Failed to download bb_core.yml", nil, "status", resp.StatusCode, "url", bbCoreURL)
		return fmt.Errorf("failed to download bb_core.yml: HTTP %d", resp.StatusCode)
	}

	file, err := os.Create(bbCorePath)
	if err != nil {
		utils.LogError("Failed to create bb_core.yml file", err, "path", bbCorePath)
		return fmt.Errorf("failed to create bb_core.yml file: %v", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		utils.LogError("Failed to write bb_core.yml file", err, "path", bbCorePath)
		return fmt.Errorf("failed to write bb_core.yml file: %v", err)
	}

	utils.LogInfo("Core variables installed successfully online", "path", bbCorePath)
	fmt.Println("Core variables installed successfully.")
	return nil
}

func InstallCoreVariablesOffline(coreVarsPath, userHome string) error {
	utils.LogInfo("Installing core variables offline", "core_vars_path", coreVarsPath, "home", userHome)

	// Validate userHome is not empty
	if userHome == "" {
		utils.LogError("User home directory is empty", nil)
		return fmt.Errorf("user home directory cannot be empty")
	}

	// Create inventory directory structure
	inventoryDir := filepath.Join(userHome, "bluebanquise", "inventory")
	groupVarsDir := filepath.Join(inventoryDir, "group_vars", "all")

	utils.LogInfo("Creating inventory directory structure", "path", groupVarsDir)
	if err := os.MkdirAll(groupVarsDir, 0755); err != nil {
		utils.LogError("Failed to create inventory directory", err, "path", groupVarsDir)
		return fmt.Errorf("failed to create inventory directory: %v", err)
	}

	// Check if coreVarsPath is a file or directory
	info, err := os.Stat(coreVarsPath)
	if err != nil {
		utils.LogError("Failed to stat core variables path", err, "path", coreVarsPath)
		return fmt.Errorf("failed to stat core variables path: %v", err)
	}

	if info.IsDir() {
		// Directory containing multiple variable files
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
		// Single variable file
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
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
