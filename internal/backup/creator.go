package backup

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// CreateBackup creates a zip backup of the profile directory
func CreateBackup(profilePath, profileName string) (string, error) {
	// Create backup directory if it doesn't exist
	backupDir := "backup"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102-1504")
	backupFilename := fmt.Sprintf("settings_%s_%s.zip", profileName, timestamp)
	backupPath := filepath.Join(backupDir, backupFilename)

	// Create zip file
	zipFile, err := os.Create(backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk through profile directory and add files to zip
	err = filepath.Walk(profilePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Calculate relative path for zip
		relPath, err := filepath.Rel(profilePath, path)
		if err != nil {
			return err
		}

		// Create file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = relPath
		header.Method = zip.Deflate

		// Create writer for file
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Open source file
		sourceFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		// Copy file content to zip
		_, err = io.Copy(writer, sourceFile)
		return err
	})

	if err != nil {
		return "", fmt.Errorf("failed to create backup: %w", err)
	}

	// Close zip writer to finalize
	if err := zipWriter.Close(); err != nil {
		return "", fmt.Errorf("failed to finalize backup: %w", err)
	}

	return backupPath, nil
}

// VerifyBackup verifies that a backup file exists and is readable
func VerifyBackup(backupPath string) error {
	// Check if file exists
	info, err := os.Stat(backupPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("backup file does not exist: %s", backupPath)
		}
		return fmt.Errorf("failed to access backup file: %w", err)
	}

	// Check if file is not empty
	if info.Size() == 0 {
		return fmt.Errorf("backup file is empty: %s", backupPath)
	}

	// Try to open and read zip file
	zipReader, err := zip.OpenReader(backupPath)
	if err != nil {
		return fmt.Errorf("backup file is not a valid zip: %w", err)
	}
	defer zipReader.Close()

	// Check if zip has at least one file
	if len(zipReader.File) == 0 {
		return fmt.Errorf("backup zip file is empty")
	}

	return nil
}
