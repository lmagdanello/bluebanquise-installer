package cmd

import (
	"fmt"
	"os"

	"github.com/lmagdanello/bluebanquise-installer/internal/bootstrap"
	"github.com/lmagdanello/bluebanquise-installer/internal/system"
	"github.com/lmagdanello/bluebanquise-installer/internal/utils"
	"github.com/spf13/cobra"
)

var (
	collectionsPath        string
	tarballPath            string
	requirementsPath       string
	coreVarsPath           string
	userName               string
	userHome               string
	offlineSkipEnvironment bool
	offlineDebug           bool
)

var offlineCmd = &cobra.Command{
	Use:   "offline",
	Short: "Install BlueBanquise in offline mode",
	Long: `Install BlueBanquise in offline mode using local collections, tarballs, and requirements.
	
This command will:
1. Check system prerequisites
2. Validate collections path or tarball path
3. Validate requirements path (if provided)
4. Detect the operating system
5. Install required system packages
6. Create bluebanquise user
7. Configure Python virtual environment (with offline requirements if provided)
8. Install BlueBanquise collections from local path or tarballs

You can use either --collections-path (for pre-installed collections) or --tarball-path (for tarball files).
You can optionally use --requirements-path for offline Python package installation.`,
	Run: func(cmd *cobra.Command, args []string) {
		if collectionsPath == "" && tarballPath == "" {
			utils.LogError("Missing required path", nil, "collections_path", collectionsPath, "tarball_path", tarballPath)
			fmt.Println("Error: either --collections-path or --tarball-path is required for offline installation")
			os.Exit(1)
		}

		if collectionsPath != "" && tarballPath != "" {
			utils.LogError("Both paths provided", nil, "collections_path", collectionsPath, "tarball_path", tarballPath)
			fmt.Println("Error: cannot use both --collections-path and --tarball-path at the same time")
			os.Exit(1)
		}

		utils.LogInfo("Starting BlueBanquise offline installation",
			"collections_path", collectionsPath,
			"tarball_path", tarballPath,
			"requirements_path", requirementsPath,
			"user", userName,
			"home", userHome,
			"skip_environment", offlineSkipEnvironment,
			"debug", offlineDebug)

		// Check system prerequisites (without internet connectivity)
		utils.LogInfo("Checking system prerequisites")
		fmt.Println("Checking system prerequisites...")
		if err := utils.SystemCheck(); err != nil {
			utils.LogError("System check failed", err)
			fmt.Printf("System check failed: %v\n", err)
			os.Exit(1)
		}

		// Validate input path
		if tarballPath != "" {
			utils.LogInfo("Validating tarball path", "path", tarballPath)
			fmt.Println("Validating tarball path...")
			if err := utils.CheckTarballPrerequisites(tarballPath); err != nil {
				utils.LogError("Tarball validation failed", err, "path", tarballPath)
				fmt.Printf("Tarball validation failed: %v\n", err)
				os.Exit(1)
			}
		} else {
			utils.LogInfo("Validating collections path", "path", collectionsPath)
			fmt.Println("Validating collections path...")
			if err := utils.CheckOfflinePrerequisites(collectionsPath); err != nil {
				utils.LogError("Collections validation failed", err, "path", collectionsPath)
				fmt.Printf("Collections validation failed: %v\n", err)
				os.Exit(1)
			}
		}

		// Validate requirements path if provided
		if requirementsPath != "" {
			utils.LogInfo("Validating requirements path", "path", requirementsPath)
			fmt.Println("Validating requirements path...")
			if err := utils.CheckRequirementsPrerequisites(requirementsPath); err != nil {
				utils.LogError("Requirements validation failed", err, "path", requirementsPath)
				fmt.Printf("Requirements validation failed: %v\n", err)
				os.Exit(1)
			}
		}

		// Validate core variables path if provided
		if coreVarsPath != "" {
			utils.LogInfo("Validating core variables path", "path", coreVarsPath)
			fmt.Println("Validating core variables path...")
			if _, err := os.Stat(coreVarsPath); err != nil {
				utils.LogError("Core variables path validation failed", err, "path", coreVarsPath)
				fmt.Printf("Core variables path validation failed: %v\n", err)
				os.Exit(1)
			}
		}

		// Detect OS
		utils.LogInfo("Detecting operating system")
		osID, version, err := system.DetectOS()
		if err != nil {
			utils.LogError("Error detecting OS", err)
			fmt.Printf("Error detecting OS: %v\n", err)
			os.Exit(1)
		}
		utils.LogInfo("OS detected", "os", osID, "version", version)
		fmt.Printf("Detected OS: %s %s\n", osID, version)

		// Find packages for this OS
		var packages []string
		for _, pkg := range system.DependenciePackages {
			if pkg.OSID == osID && pkg.Version == version {
				packages = pkg.Packages
				break
			}
		}

		if len(packages) == 0 {
			utils.LogError("No package definition found", nil, "os", osID, "version", version)
			fmt.Printf("No package definition found for %s %s\n", osID, version)
			os.Exit(1)
		}

		// Install system packages
		utils.LogInfo("Installing system packages", "packages", packages)
		fmt.Println("Installing system packages...")
		if err := utils.InstallPackages(packages); err != nil {
			utils.LogError("Error installing packages", err, "packages", packages)
			fmt.Printf("Error installing packages: %v\n", err)
			os.Exit(1)
		}

		// Create bluebanquise user
		utils.LogInfo("Creating BlueBanquise user", "user", userName, "home", userHome)
		if err := bootstrap.CreateBluebanquiseUser(userName, userHome); err != nil {
			utils.LogError("Error creating user", err, "user", userName, "home", userHome)
			fmt.Printf("Error creating user: %v\n", err)
			os.Exit(1)
		}

		// Configure environment (unless skipped)
		if !offlineSkipEnvironment {
			utils.LogInfo("Configuring environment")
			if err := bootstrap.ConfigureEnvironmentOffline(userName, userHome, requirementsPath); err != nil {
				utils.LogError("Error configuring environment", err)
				fmt.Printf("Error configuring environment: %v\n", err)
				os.Exit(1)
			}
		} else {
			utils.LogInfo("Skipping environment configuration")
		}

		// Install collections (requires environment to be configured)
		if tarballPath != "" {
			utils.LogInfo("Installing collections from tarballs", "path", tarballPath)
			if err := bootstrap.InstallCollectionsFromPath(tarballPath, userHome); err != nil {
				utils.LogError("Error installing collections from tarballs", err, "path", tarballPath)
				fmt.Printf("Error installing collections from tarballs: %v\n", err)
				os.Exit(1)
			}
		} else {
			utils.LogInfo("Installing collections offline", "path", collectionsPath)
			if err := bootstrap.InstallCollectionsFromPath(collectionsPath, userHome); err != nil {
				utils.LogError("Error installing collections", err, "path", collectionsPath)
				fmt.Printf("Error installing collections: %v\n", err)
				os.Exit(1)
			}
		}

		// Install core variables offline if path provided
		if coreVarsPath != "" {
			utils.LogInfo("Installing core variables offline")
			if err := bootstrap.InstallCoreVariablesOffline(coreVarsPath, userHome); err != nil {
				utils.LogError("Error installing core variables", err)
				fmt.Printf("Error installing core variables: %v\n", err)
				os.Exit(1)
			}
		} else {
			utils.LogInfo("No core variables path provided, skipping core variables installation")
		}

		utils.LogInfo("Offline installation completed successfully")
		utils.ShowCompletionMessage(userName, userHome)
	},
}

func init() {
	offlineCmd.Flags().StringVarP(&collectionsPath, "collections-path", "c", "", "Path to local BlueBanquise collections")
	offlineCmd.Flags().StringVarP(&tarballPath, "tarball-path", "t", "", "Path to BlueBanquise collection tarballs")
	offlineCmd.Flags().StringVarP(&requirementsPath, "requirements-path", "r", "", "Path to Python requirements for offline installation")
	offlineCmd.Flags().StringVarP(&coreVarsPath, "core-vars-path", "v", "", "Path to core variables for offline installation")
	offlineCmd.Flags().StringVarP(&userName, "user", "u", "bluebanquise", "Username for BlueBanquise")
	offlineCmd.Flags().StringVarP(&userHome, "home", "H", "/var/lib/bluebanquise", "Home directory for BlueBanquise user")
	offlineCmd.Flags().BoolVarP(&offlineSkipEnvironment, "skip-environment", "e", false, "Skip environment configuration")
	offlineCmd.Flags().BoolVarP(&offlineDebug, "debug", "d", false, "Enable debug mode")

	rootCmd.AddCommand(offlineCmd)
}
