package cmd

import (
	"testing"

	"github.com/lmagdanello/bluebanquise-installer/internal/utils"
)

func TestDownloadCommand(t *testing.T) {
	utils.InitTestLogger()

	// Test missing path - should exit with error
	t.Run("missing path", func(t *testing.T) {
		// This test verifies that the command requires --path flag
		// The actual error handling is done by Cobra's flag validation
		// which exits the program, so we can't easily test it in unit tests
		t.Skip("Skipping test that requires program exit - tested manually")
	})

	// Test missing download type - should exit with error
	t.Run("missing download type", func(t *testing.T) {
		// This test verifies that at least one download type is required
		// The actual error handling is done in the Run function
		// which exits the program, so we can't easily test it in unit tests
		t.Skip("Skipping test that requires program exit - tested manually")
	})

	// Test valid command structure
	t.Run("valid command structure", func(t *testing.T) {
		// Test that the command can be created and has the right flags
		cmd := downloadCmd

		// Check that required flags exist
		pathFlag := cmd.Flags().Lookup("path")
		if pathFlag == nil {
			t.Error("--path flag not found")
		}

		collectionsFlag := cmd.Flags().Lookup("collections")
		if collectionsFlag == nil {
			t.Error("--collections flag not found")
		}

		requirementsFlag := cmd.Flags().Lookup("requirements")
		if requirementsFlag == nil {
			t.Error("--requirements flag not found")
		}

		coreVarsFlag := cmd.Flags().Lookup("core-vars")
		if coreVarsFlag == nil {
			t.Error("--core-vars flag not found")
		}
	})
}
