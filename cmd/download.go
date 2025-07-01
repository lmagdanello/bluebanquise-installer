package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lmagdanello/bluebanquise-installer/internal/system"
	"github.com/lmagdanello/bluebanquise-installer/internal/utils"
	"github.com/spf13/cobra"
)

var (
	downloadPath         string
	downloadCollections  bool
	downloadRequirements bool
	downloadCoreVars     bool
	downloadCmd          = &cobra.Command{
		Use:   "download",
		Short: "Download BlueBanquise collections and requirements for offline installation",
		Long: `Download BlueBanquise collections and requirements from GitHub for offline installation.
		
This command downloads files to a base directory specified by --path. Use specific flags to download different components:

--collections: Downloads collection tarballs to <path>/collections/
--requirements: Downloads Python packages to <path>/requirements/
--core-vars: Downloads core variables to <path>/core-vars/

You can use multiple flags to download multiple components at once.

Examples:
  # Download collections only
  ./bluebanquise-installer download --path /tmp/offline --collections

  # Download requirements only
  ./bluebanquise-installer download --path /tmp/requirements --requirements

  # Download core variables only
  ./bluebanquise-installer download --path /tmp/core-vars --core-vars

  # Download everything
  ./bluebanquise-installer download --path /tmp/offline --collections --requirements --core-vars`,
		Run: func(cmd *cobra.Command, args []string) {
			if downloadPath == "" {
				utils.LogError("Missing download path", nil)
				fmt.Println("Error: --path is required")
				os.Exit(1)
			}

			if !downloadCollections && !downloadRequirements && !downloadCoreVars {
				utils.LogError("No download type specified", nil)
				fmt.Println("Error: specify at least one of --collections, --requirements, or --core-vars")
				os.Exit(1)
			}

			utils.LogInfo("Starting BlueBanquise download",
				"path", downloadPath,
				"collections", downloadCollections,
				"requirements", downloadRequirements,
				"core-vars", downloadCoreVars)

			// Create base download directory
			if err := os.MkdirAll(downloadPath, 0755); err != nil {
				utils.LogError("Error creating download directory", err, "path", downloadPath)
				fmt.Printf("Error creating download directory: %v\n", err)
				os.Exit(1)
			}

			if downloadCollections {
				downloadCollectionsToPath()
			}
			if downloadRequirements {
				downloadRequirementsToPath()
			}
			if downloadCoreVars {
				downloadCoreVarsToPath()
			}
		},
	}
)

func downloadCollectionsToPath() {
	collectionsPath := filepath.Join(downloadPath, "collections")
	utils.LogInfo("Downloading collections", "path", collectionsPath)

	// Create temporary Python environment outside download directory
	tempVenv := filepath.Join(os.TempDir(), "bluebanquise_download_venv")
	if err := utils.RunCommand("/usr/bin/python3", "-m", "venv", tempVenv); err != nil {
		utils.LogError("Error creating temporary virtual environment", err, "path", tempVenv)
		fmt.Printf("Error creating temporary virtual environment: %v\n", err)
		os.Exit(1)
	}

	// Install ansible-galaxy in temp environment
	python3 := filepath.Join(tempVenv, "bin", "python3")
	if err := utils.RunCommand(python3, "-m", "pip", "install", "ansible-core"); err != nil {
		utils.LogError("Error installing ansible-core", err)
		fmt.Printf("Error installing ansible-core: %v\n", err)
		os.Exit(1)
	}

	// Download tarballs
	ansibleGalaxy := filepath.Join(tempVenv, "bin", "ansible-galaxy")

	utils.LogInfo("Downloading BlueBanquise collection tarball")
	fmt.Println("Downloading BlueBanquise collection tarball...")
	if err := utils.RunCommand(ansibleGalaxy,
		"collection", "download",
		"git+https://github.com/bluebanquise/bluebanquise.git#/collections/infrastructure,master",
		"-p", collectionsPath); err != nil {
		utils.LogError("Error downloading BlueBanquise tarball", err)
		fmt.Printf("Error downloading BlueBanquise tarball: %v\n", err)
		os.Exit(1)
	}

	utils.LogInfo("Downloading community.general collection tarball")
	fmt.Println("Downloading community.general collection tarball...")
	if err := utils.RunCommand(ansibleGalaxy,
		"collection", "download",
		"community.general",
		"-p", collectionsPath); err != nil {
		utils.LogError("Error downloading community.general tarball", err)
		fmt.Printf("Error downloading community.general tarball: %v\n", err)
		os.Exit(1)
	}

	// Clean up temp environment
	if err := os.RemoveAll(tempVenv); err != nil {
		utils.LogWarning("Could not remove temporary environment", "error", err, "path", tempVenv)
		fmt.Printf("Warning: could not remove temporary environment: %v\n", err)
	}

	utils.LogInfo("Collections downloaded successfully", "path", collectionsPath)
	fmt.Printf("Collections downloaded successfully to: %s\n", collectionsPath)
	fmt.Println("Transfer this directory to your target machine and use with:")
	fmt.Printf("  ./bluebanquise-installer offline --collections-path %s\n", collectionsPath)
}

func downloadRequirementsToPath() {
	requirementsPath := filepath.Join(downloadPath, "requirements")
	utils.LogInfo("Downloading Python requirements", "path", requirementsPath)

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

	if err := utils.DownloadRequirements(requirements, requirementsPath); err != nil {
		utils.LogError("Error downloading requirements", err)
		fmt.Printf("Error downloading requirements: %v\n", err)
		os.Exit(1)
	}

	utils.LogInfo("Python requirements downloaded successfully", "path", requirementsPath)
	fmt.Printf("Python requirements downloaded successfully to: %s\n", requirementsPath)
	fmt.Println("Transfer this directory to your target machine and use with:")
	fmt.Printf("  ./bluebanquise-installer offline --collections-path <collections-path> --requirements-path %s\n", requirementsPath)
}

func downloadCoreVarsToPath() {
	coreVarsPath := filepath.Join(downloadPath, "core-vars")
	utils.LogInfo("Downloading core variables", "path", coreVarsPath)

	// Download core variables from GitHub
	utils.LogInfo("Downloading core variables from GitHub")
	fmt.Println("Downloading core variables from GitHub...")
	if err := utils.DownloadFile("https://raw.githubusercontent.com/bluebanquise/bluebanquise/refs/heads/master/resources/bb_core.yml", filepath.Join(coreVarsPath, "bb_core.yml")); err != nil {
		utils.LogError("Error downloading core variables", err)
		fmt.Printf("Error downloading core variables: %v\n", err)
		os.Exit(1)
	}

	utils.LogInfo("Core variables downloaded successfully", "path", coreVarsPath)
	fmt.Printf("Core variables downloaded successfully to: %s\n", coreVarsPath)
	fmt.Println("Transfer this file to your target machine and use with:")
	fmt.Printf("  ./bluebanquise-installer offline --collections-path <collections-path> --core-vars-path %s/bb_core.yml\n", coreVarsPath)
}

func init() {
	downloadCmd.Flags().StringVarP(&downloadPath, "path", "p", "", "Path to download collections (required)")
	downloadCmd.Flags().BoolVarP(&downloadCollections, "collections", "c", false, "Download collections/tarballs for offline installation")
	downloadCmd.Flags().BoolVarP(&downloadRequirements, "requirements", "r", false, "Download Python requirements for offline installation")
	downloadCmd.Flags().BoolVarP(&downloadCoreVars, "core-vars", "v", false, "Download core variables for offline installation")
	if err := downloadCmd.MarkFlagRequired("path"); err != nil {
		utils.LogError("Error marking path flag as required", err)
		os.Exit(1)
	}

	rootCmd.AddCommand(downloadCmd)
}
