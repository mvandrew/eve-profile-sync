package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// UserFile represents a user configuration file
type UserFile struct {
	ID   string
	Name string // May be empty if extraction fails
	Path string
}

// CharacterFile represents a character configuration file
type CharacterFile struct {
	ID   string
	Name string // May be empty if extraction fails
	Path string
}

// ListUserFiles lists all user files in a profile directory
func ListUserFiles(profilePath string) ([]UserFile, error) {
	if err := ValidateProfilesDirectory(profilePath); err != nil {
		return nil, err
	}

	var userFiles []UserFile

	entries, err := os.ReadDir(profilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Check if file matches pattern core_user_*.dat
		if strings.HasPrefix(name, "core_user_") && strings.HasSuffix(name, ".dat") {
			userID, err := ExtractUserID(name)
			if err != nil {
				// Skip invalid files
				continue
			}

			filePath := filepath.Join(profilePath, name)

			// Try to extract user name (may fail for encrypted files)
			userName, _ := TryExtractUserName(filePath)

			userFiles = append(userFiles, UserFile{
				ID:   userID,
				Name: userName,
				Path: filePath,
			})
		}
	}

	return userFiles, nil
}

// ListCharacterFiles lists all character files in a profile directory
func ListCharacterFiles(profilePath string) ([]CharacterFile, error) {
	if err := ValidateProfilesDirectory(profilePath); err != nil {
		return nil, err
	}

	var charFiles []CharacterFile

	entries, err := os.ReadDir(profilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Check if file matches pattern core_char_*.dat
		if strings.HasPrefix(name, "core_char_") && strings.HasSuffix(name, ".dat") {
			charID, err := ExtractCharacterID(name)
			if err != nil {
				// Skip invalid files
				continue
			}

			filePath := filepath.Join(profilePath, name)

			// Try to extract character name (may fail for encrypted files)
			charName, _ := TryExtractCharacterName(filePath)

			charFiles = append(charFiles, CharacterFile{
				ID:   charID,
				Name: charName,
				Path: filePath,
			})
		}
	}

	return charFiles, nil
}
