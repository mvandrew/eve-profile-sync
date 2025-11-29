package sync

import (
	"fmt"
	"os"
	"path/filepath"
)

// ValidateOperation validates that all prerequisites for sync operation are met
func ValidateOperation(profilePath, sourceUserFile, sourceCharFile string) error {
	// Validate profile path
	if err := ValidateProfilePath(profilePath); err != nil {
		return fmt.Errorf("profile validation failed: %w", err)
	}

	// Validate source user file
	if err := ValidateSourceFile(sourceUserFile); err != nil {
		return fmt.Errorf("user file validation failed: %w", err)
	}

	// Validate source character file
	if err := ValidateSourceFile(sourceCharFile); err != nil {
		return fmt.Errorf("character file validation failed: %w", err)
	}

	// Check if profile directory is writable
	if err := checkWritable(profilePath); err != nil {
		return fmt.Errorf("profile directory is not writable: %w", err)
	}

	return nil
}

// checkWritable checks if a directory is writable
func checkWritable(dirPath string) error {
	testFile := filepath.Join(dirPath, ".write_test")

	// Try to create a test file
	file, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("cannot write to directory: %w", err)
	}
	file.Close()

	// Clean up test file
	if err := os.Remove(testFile); err != nil {
		// Log but don't fail on cleanup error
		_ = err
	}

	return nil
}
