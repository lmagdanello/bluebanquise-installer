package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/lmagdanello/bluebanquise-installer/internal/system"
	"github.com/lmagdanello/bluebanquise-installer/internal/utils"
	"github.com/spf13/cobra"
)

const (
	nonExistentPath = "/non/existent/path"
)

var (
	downloadPath         string
	downloadTarball      bool
	downloadRequirements bool
	downloadCoreVars     bool
	downloadCmd          = &cobra.Command{
		Use:   "download",
		Short: "Download BlueBanquise collections and requirements for offline installation",
		Long: `Download BlueBanquise collections and requirements from GitHub for offline installation.
		
This command provides four methods:

Method 1 - Download collections to directory (default):
  Downloads collections using ansible-galaxy and prepares them for offline installation.

Method 2 - Download tarballs:
  Downloads collection tarballs that can be manually transferred and installed.
  Use --tarball flag for this method.

Method 3 - Download Python requirements:
  Downloads Python packages required by BlueBanquise for offline installation.
  Use --requirements flag for this method.

Method 4 - Download core variables:
  Downloads core variables (bb_core.yml) from GitHub for offline installation.
  Use --core-vars flag for this method.

Examples:
  # Download collections to directory
  ./bluebanquise-installer download --path /tmp/bluebanquise-offline

  # Download tarballs
  ./bluebanquise-installer download --path /tmp/tarballs --tarball

  # Download Python requirements
  ./bluebanquise-installer download --path /tmp/requirements --requirements

  # Download core variables
  ./bluebanquise-installer download --path /tmp/core-vars --core-vars`,
		Run: func(cmd *cobra.Command, args []string) {
			if downloadPath == "" {
				utils.LogError("Missing download path", nil)
				fmt.Println("Error: --path is required")
				os.Exit(1)
			}

			utils.LogInfo("Starting BlueBanquise download",
				"path", downloadPath,
				"tarball", downloadTarball,
				"requirements", downloadRequirements,
				"core-vars", downloadCoreVars)

			// Create download directory
			if err := os.MkdirAll(downloadPath, 0755); err != nil {
				utils.LogError("Error creating download directory", err, "path", downloadPath)
				fmt.Printf("Error creating download directory: %v\n", err)
				os.Exit(1)
			}

			if downloadRequirements {
				downloadPythonRequirements()
			} else if downloadTarball {
				downloadTarballs()
			} else if downloadCoreVars {
				downloadCoreVariables()
			} else {
				downloadCollections()
			}
		},
	}
)

func downloadPythonRequirements() {
	utils.LogInfo("Downloading Python requirements", "path", downloadPath)

	// Detect OS to get the correct requirements
	osID, version, err := system.DetectOS()
	if err != nil {
		utils.LogError("Error detecting OS", err)
		fmt.Printf("Error detecting OS: %v\n", err)
		os.Exit(1)
	}

	// Get requirements for this OS
	var requirements []string
	for _, pkg := range system.DependenciePackages {
		if pkg.OSID == osID && pkg.Version == version {
			requirements = system.PythonRequirements
			break
		}
	}

	if len(requirements) == 0 {
		utils.LogError("No requirements found for OS", nil, "os", osID, "version", version)
		fmt.Printf("No requirements found for %s %s\n", osID, version)
		os.Exit(1)
	}

	utils.LogInfo("Downloading requirements for OS", "os", osID, "version", version, "requirements", requirements)
	fmt.Printf("Downloading Python requirements for %s %s...\n", osID, version)

	if err := utils.DownloadRequirements(requirements, downloadPath); err != nil {
		utils.LogError("Error downloading requirements", err)
		fmt.Printf("Error downloading requirements: %v\n", err)
		os.Exit(1)
	}

	utils.LogInfo("Python requirements downloaded successfully", "path", downloadPath)
	fmt.Printf("Python requirements downloaded successfully to: %s\n", downloadPath)
	fmt.Println("Transfer this directory to your target machine and use with:")
	fmt.Printf("  ./bluebanquise-installer offline --collections-path <collections-path> --requirements-path %s\n", downloadPath)
}

func downloadCollections() {
	utils.LogInfo("Downloading collections to directory", "path", downloadPath)

	// Create temporary Python environment
	tempVenv := filepath.Join(downloadPath, "temp_venv")
	if err := utils.RunCommand("/usr/bin/python3", "-m", "venv", tempVenv); err != nil {
		utils.LogError("Error creating temporary virtual environment", err, "path", tempVenv)
		fmt.Printf("Error creating temporary virtual environment: %v\n", err)
		os.Exit(1)
	}

	// Install ansible-galaxy in temp environment
	pip := filepath.Join(tempVenv, "bin", "pip")
	if err := utils.RunCommand(pip, "install", "ansible-core"); err != nil {
		utils.LogError("Error installing ansible-core", err)
		fmt.Printf("Error installing ansible-core: %v\n", err)
		os.Exit(1)
	}

	// Download collections
	ansibleGalaxy := filepath.Join(tempVenv, "bin", "ansible-galaxy")
	collectionsDir := filepath.Join(downloadPath, "collections")

	utils.LogInfo("Downloading BlueBanquise collections")
	fmt.Println("Downloading BlueBanquise collections...")
	downloadCmd := exec.Command(ansibleGalaxy,
		"collection", "download",
		"git+https://github.com/bluebanquise/bluebanquise.git#/collections/infrastructure,master",
		"-p", collectionsDir)
	if err := downloadCmd.Run(); err != nil {
		utils.LogError("Error downloading collections", err)
		fmt.Printf("Error downloading collections: %v\n", err)
		os.Exit(1)
	}

	// Download community.general
	utils.LogInfo("Downloading community.general collection")
	fmt.Println("Downloading community.general collection...")
	communityCmd := exec.Command(ansibleGalaxy,
		"collection", "download",
		"community.general",
		"-p", collectionsDir)
	if err := communityCmd.Run(); err != nil {
		utils.LogError("Error downloading community.general", err)
		fmt.Printf("Error downloading community.general: %v\n", err)
		os.Exit(1)
	}

	// Clean up temp environment
	if err := os.RemoveAll(tempVenv); err != nil {
		utils.LogWarning("Could not remove temporary environment", "error", err, "path", tempVenv)
		fmt.Printf("Warning: could not remove temporary environment: %v\n", err)
	}

	utils.LogInfo("Collections downloaded successfully", "path", collectionsDir)
	fmt.Printf("Collections downloaded successfully to: %s\n", collectionsDir)
	fmt.Println("You can now use this path with the offline command:")
	fmt.Printf("  ./bluebanquise-installer offline --collections-path %s\n", collectionsDir)
}

func downloadTarballs() {
	utils.LogInfo("Downloading tarballs", "path", downloadPath)

	// Create temporary Python environment
	tempVenv := filepath.Join(downloadPath, "temp_venv")
	if err := utils.RunCommand("/usr/bin/python3", "-m", "venv", tempVenv); err != nil {
		utils.LogError("Error creating temporary virtual environment", err, "path", tempVenv)
		fmt.Printf("Error creating temporary virtual environment: %v\n", err)
		os.Exit(1)
	}

	// Install ansible-galaxy in temp environment
	pip := filepath.Join(tempVenv, "bin", "pip")
	if err := utils.RunCommand(pip, "install", "ansible-core"); err != nil {
		utils.LogError("Error installing ansible-core", err)
		fmt.Printf("Error installing ansible-core: %v\n", err)
		os.Exit(1)
	}

	// Download tarballs
	ansibleGalaxy := filepath.Join(tempVenv, "bin", "ansible-galaxy")

	utils.LogInfo("Downloading BlueBanquise collection tarball")
	fmt.Println("Downloading BlueBanquise collection tarball...")
	bluebanquiseCmd := exec.Command(ansibleGalaxy,
		"collection", "download",
		"git+https://github.com/bluebanquise/bluebanquise.git#/collections/infrastructure,master",
		"-p", downloadPath)
	if err := bluebanquiseCmd.Run(); err != nil {
		utils.LogError("Error downloading BlueBanquise tarball", err)
		fmt.Printf("Error downloading BlueBanquise tarball: %v\n", err)
		os.Exit(1)
	}

	utils.LogInfo("Downloading community.general collection tarball")
	fmt.Println("Downloading community.general collection tarball...")
	communityCmd := exec.Command(ansibleGalaxy,
		"collection", "download",
		"community.general",
		"-p", downloadPath)
	if err := communityCmd.Run(); err != nil {
		utils.LogError("Error downloading community.general tarball", err)
		fmt.Printf("Error downloading community.general tarball: %v\n", err)
		os.Exit(1)
	}

	// Clean up temp environment
	if err := os.RemoveAll(tempVenv); err != nil {
		utils.LogWarning("Could not remove temporary environment", "error", err, "path", tempVenv)
		fmt.Printf("Warning: could not remove temporary environment: %v\n", err)
	}

	utils.LogInfo("Tarballs downloaded successfully", "path", downloadPath)
	fmt.Printf("Tarballs downloaded successfully to: %s\n", downloadPath)
	fmt.Println("Transfer these files to your target machine and install with:")
	fmt.Println("  ansible-galaxy collection install <tarball-file> -p <collections-directory>")
	fmt.Println("  ./bluebanquise-installer offline --collections-path <collections-directory>")
}

func downloadCoreVariables() {
	utils.LogInfo("Downloading core variables", "path", downloadPath)

	// Download core variables from GitHub
	utils.LogInfo("Downloading core variables from GitHub")
	fmt.Println("Downloading core variables from GitHub...")
	if err := utils.DownloadFile("https://raw.githubusercontent.com/bluebanquise/bluebanquise/refs/heads/master/resources/bb_core.yml", filepath.Join(downloadPath, "bb_core.yml")); err != nil {
		utils.LogError("Error downloading core variables", err)
		fmt.Printf("Error downloading core variables: %v\n", err)
		os.Exit(1)
	}

	utils.LogInfo("Core variables downloaded successfully", "path", downloadPath)
	fmt.Printf("Core variables downloaded successfully to: %s\n", downloadPath)
	fmt.Println("Transfer this file to your target machine and use with:")
	fmt.Printf("  ./bluebanquise-installer offline --core-vars-path %s\n", downloadPath)
}

func init() {
	downloadCmd.Flags().StringVarP(&downloadPath, "path", "p", "", "Path to download collections (required)")
	downloadCmd.Flags().BoolVarP(&downloadTarball, "tarball", "t", false, "Download tarballs instead of collections directory")
	downloadCmd.Flags().BoolVarP(&downloadRequirements, "requirements", "r", false, "Download Python requirements for offline installation")
	downloadCmd.Flags().BoolVarP(&downloadCoreVars, "core-vars", "c", false, "Download core variables for offline installation")
	if err := downloadCmd.MarkFlagRequired("path"); err != nil {
		utils.LogError("Error marking path flag as required", err)
		os.Exit(1)
	}

	rootCmd.AddCommand(downloadCmd)
}
