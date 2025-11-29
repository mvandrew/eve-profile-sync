package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Profile represents an EVE Online profile
type Profile struct {
	Name string
	Path string
}

// FindProfilesDirectory attempts to find the EVE profiles directory
func FindProfilesDirectory() (string, error) {
	// Default path: C:\Users\{user}\AppData\Local\CCP\EVE\c_ccp_eve_online_tq_tranquility
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	defaultPath := filepath.Join(userHome, "AppData", "Local", "CCP", "EVE", "c_ccp_eve_online_tq_tranquility")

	// Check if default path exists
	if _, err := os.Stat(defaultPath); err == nil {
		return defaultPath, nil
	}

	// Default path not found, return empty string
	// The caller should prompt user for path
	return "", fmt.Errorf("default EVE profiles directory not found at: %s", defaultPath)
}

// ValidateProfilesDirectory checks if a directory exists and is accessible
func ValidateProfilesDirectory(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", path)
		}
		return fmt.Errorf("failed to access directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", path)
	}

	return nil
}

// ListProfiles lists all profiles in the profiles directory
func ListProfiles(profilesDir string) ([]Profile, error) {
	if err := ValidateProfilesDirectory(profilesDir); err != nil {
		return nil, err
	}

	var profiles []Profile

	entries, err := os.ReadDir(profilesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read profiles directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Check if directory matches pattern settings_*
		if strings.HasPrefix(name, "settings_") {
			profileName := strings.TrimPrefix(name, "settings_")
			profilePath := filepath.Join(profilesDir, name)
			profiles = append(profiles, Profile{
				Name: profileName,
				Path: profilePath,
			})
		}
	}

	return profiles, nil
}
