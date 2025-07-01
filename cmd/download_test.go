package cmd

import (
	"testing"
)

func TestDownloadCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedError  bool
		expectedOutput string
	}{
		{
			name:          "missing path",
			args:          []string{"download"},
			expectedError: true,
		},
		{
			name:          "missing download type",
			args:          []string{"download", "--path", "/tmp/test"},
			expectedError: true,
		},
		{
			name:          "download collections",
			args:          []string{"download", "--path", "/tmp/test", "--collections"},
			expectedError: false,
		},
		{
			name:          "download requirements",
			args:          []string{"download", "--path", "/tmp/test", "--requirements"},
			expectedError: false,
		},
		{
			name:          "download core vars",
			args:          []string{"download", "--path", "/tmp/test", "--core-vars"},
			expectedError: false,
		},
		{
			name:          "download all",
			args:          []string{"download", "--path", "/tmp/test", "--collections", "--requirements", "--core-vars"},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags for each test
			downloadPath = ""
			downloadCollections = false
			downloadRequirements = false
			downloadCoreVars = false

			// Set up command
			cmd := downloadCmd
			cmd.SetArgs(tt.args)

			// Execute command
			err := cmd.Execute()

			// Check results
			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
