package system

import (
	"log/slog"
	"os"
	"strings"
)

// OSMapping maps OS IDs to BlueBanquise compatible names.
var OSMapping = map[string]string{
	"rhel":          "rhel",
	"centos":        "rhel",
	"rocky":         "rhel",
	"almalinux":     "rhel",
	"ubuntu":        "ubuntu",
	"debian":        "debian",
	"opensuse-leap": "opensuse-leap",
	"sles":          "opensuse-leap",
}

func DetectOS() (string, string, error) {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		slog.Error("Error detecting OS", "error", err)
		return "", "", err
	}

	var name, version string
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			name = strings.TrimPrefix(line, "ID=")
			name = strings.Trim(name, "\"")
		} else if strings.HasPrefix(line, "VERSION_ID=") {
			version = strings.TrimPrefix(line, "VERSION_ID=")
			version = strings.Trim(version, "\"")
		}
	}

	// Map OS ID to BlueBanquise compatible name
	if mappedName, exists := OSMapping[name]; exists {
		name = mappedName
	}

	// Handle version mapping for RHEL derivatives
	if name == "rhel" {
		switch {
		case strings.Contains(strings.ToLower(version), "stream"):
			// Extract version number from stream
			if strings.Contains(version, "8") {
				version = "8"
			} else if strings.Contains(version, "9") {
				version = "9"
			}
		default:
			// Extract major version (e.g., 9.3 -> 9)
			if strings.Contains(version, ".") {
				version = strings.Split(version, ".")[0]
			}
		}
	}

	return name, version, nil
}
