package sync

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ReplaceUserFiles replaces all user files with the selected user file content
func ReplaceUserFiles(profilePath, sourceUserFile string) error {
	// Read source file content
	sourceContent, err := os.ReadFile(sourceUserFile)
	if err != nil {
		return fmt.Errorf("failed to read source user file: %w", err)
	}

	// Get source filename to exclude it from replacement
	sourceFilename := filepath.Base(sourceUserFile)

	// List all files in profile directory
	entries, err := os.ReadDir(profilePath)
	if err != nil {
		return fmt.Errorf("failed to read profile directory: %w", err)
	}

	var replacedCount int
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Check if file matches pattern core_user_*.dat
		if strings.HasPrefix(name, "core_user_") && strings.HasSuffix(name, ".dat") {
			// Skip source file itself
			if name == sourceFilename {
				continue
			}

			filePath := filepath.Join(profilePath, name)

			// Write source content to target file
			if err := os.WriteFile(filePath, sourceContent, 0644); err != nil {
				return fmt.Errorf("failed to replace file %s: %w", name, err)
			}

			replacedCount++
		}
	}

	if replacedCount == 0 {
		return fmt.Errorf("no user files found to replace")
	}

	return nil
}

// ReplaceCharacterFiles replaces all character files with the selected character file content
func ReplaceCharacterFiles(profilePath, sourceCharFile string) error {
	// Read source file content
	sourceContent, err := os.ReadFile(sourceCharFile)
	if err != nil {
		return fmt.Errorf("failed to read source character file: %w", err)
	}

	// Get source filename to exclude it from replacement
	sourceFilename := filepath.Base(sourceCharFile)

	// List all files in profile directory
	entries, err := os.ReadDir(profilePath)
	if err != nil {
		return fmt.Errorf("failed to read profile directory: %w", err)
	}

	var replacedCount int
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Check if file matches pattern core_char_*.dat
		if strings.HasPrefix(name, "core_char_") && strings.HasSuffix(name, ".dat") {
			// Skip source file itself
			if name == sourceFilename {
				continue
			}

			filePath := filepath.Join(profilePath, name)

			// Write source content to target file
			if err := os.WriteFile(filePath, sourceContent, 0644); err != nil {
				return fmt.Errorf("failed to replace file %s: %w", name, err)
			}

			replacedCount++
		}
	}

	if replacedCount == 0 {
		return fmt.Errorf("no character files found to replace")
	}

	return nil
}

// ValidateProfilePath validates that a profile path exists and is accessible
func ValidateProfilePath(profilePath string) error {
	info, err := os.Stat(profilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("profile path does not exist: %s", profilePath)
		}
		return fmt.Errorf("failed to access profile path: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("profile path is not a directory: %s", profilePath)
	}

	return nil
}

// ValidateSourceFile validates that a source file exists and is readable
func ValidateSourceFile(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("source file does not exist: %s", filePath)
		}
		return fmt.Errorf("failed to access source file: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("source path is a directory, not a file: %s", filePath)
	}

	// Try to read file to ensure it's accessible
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer file.Close()

	// Read a small portion to verify readability
	buf := make([]byte, 1)
	_, err = file.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	return nil
}
