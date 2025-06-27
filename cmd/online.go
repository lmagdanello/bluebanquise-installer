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
	onlineUserName        string
	onlineUserHome        string
	onlineSkipEnvironment bool
	onlineDebug           bool
)

var onlineCmd = &cobra.Command{
	Use:   "online",
	Short: "Install BlueBanquise in online mode",
	Long: `Install BlueBanquise in online mode downloading collections from GitHub.
	
	This command will:
	1. Check system prerequisites
	2. Detect the operating system
	3. Install required system packages
	4. Create bluebanquise user
	5. Configure Python virtual environment
	6. Install BlueBanquise collections from GitHub`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.LogInfo("Starting BlueBanquise online installation",
			"user", onlineUserName,
			"home", onlineUserHome,
			"skip_environment", onlineSkipEnvironment,
			"debug", onlineDebug)

		// Check system prerequisites
		utils.LogInfo("Checking system prerequisites")
		fmt.Println("Checking system prerequisites...")
		if err := utils.SystemCheck(); err != nil {
			utils.LogError("System check failed", err)
			fmt.Printf("System check failed: %v\n", err)
			os.Exit(1)
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
		var postHook func() error
		for _, pkg := range system.DependenciePackages {
			if pkg.OSID == osID && pkg.Version == version {
				packages = pkg.Packages
				postHook = pkg.PostHook
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

		// Run post-installation hook if exists
		if postHook != nil {
			utils.LogInfo("Running post-installation hook")
			fmt.Println("Running post-installation hook...")
			if err := postHook(); err != nil {
				utils.LogError("Error in post-installation hook", err)
				fmt.Printf("Error in post-installation hook: %v\n", err)
				os.Exit(1)
			}
		}

		// Create bluebanquise user
		utils.LogInfo("Creating BlueBanquise user", "user", onlineUserName, "home", onlineUserHome)
		if err := bootstrap.CreateBluebanquiseUser(onlineUserName, onlineUserHome); err != nil {
			utils.LogError("Error creating user", err, "user", onlineUserName, "home", onlineUserHome)
			fmt.Printf("Error creating user: %v\n", err)
			os.Exit(1)
		}

		// Configure environment (unless skipped)
		if !onlineSkipEnvironment {
			utils.LogInfo("Configuring environment")
			if err := bootstrap.ConfigureEnvironment(onlineUserName, onlineUserHome, ""); err != nil {
				utils.LogError("Error configuring environment", err)
				fmt.Printf("Error configuring environment: %v\n", err)
				os.Exit(1)
			}
		} else {
			utils.LogInfo("Skipping environment configuration")
		}

		// Install collections online
		utils.LogInfo("Installing collections online")
		if err := bootstrap.InstallCollectionsOnline(onlineUserHome); err != nil {
			utils.LogError("Error installing collections", err)
			fmt.Printf("Error installing collections: %v\n", err)
			os.Exit(1)
		}

		// Install core variables online
		utils.LogInfo("Installing core variables online")
		if err := bootstrap.InstallCoreVariablesOnline(onlineUserHome); err != nil {
			utils.LogError("Error installing core variables", err)
			fmt.Printf("Error installing core variables: %v\n", err)
			os.Exit(1)
		}

		utils.LogInfo("Online installation completed successfully")
		utils.ShowCompletionMessage(onlineUserName, onlineUserHome)
	},
}

func init() {
	onlineCmd.Flags().StringVarP(&onlineUserName, "user", "u", "bluebanquise", "Username for BlueBanquise")
	onlineCmd.Flags().StringVarP(&onlineUserHome, "home", "h", "/var/lib/bluebanquise", "Home directory for BlueBanquise user")
	onlineCmd.Flags().BoolVarP(&onlineSkipEnvironment, "skip-environment", "e", false, "Skip environment configuration")
	onlineCmd.Flags().BoolVarP(&onlineDebug, "debug", "d", false, "Enable debug mode")

	rootCmd.AddCommand(onlineCmd)
}
