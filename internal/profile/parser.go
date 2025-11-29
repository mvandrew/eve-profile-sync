package profile

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	userIDRegex      = regexp.MustCompile(`^core_user_(\d+)\.dat$`)
	characterIDRegex = regexp.MustCompile(`^core_char_(\d+)\.dat$`)
)

// ExtractUserID extracts user ID from filename
func ExtractUserID(filename string) (string, error) {
	matches := userIDRegex.FindStringSubmatch(filename)
	if len(matches) < 2 {
		return "", fmt.Errorf("invalid user filename format: %s", filename)
	}
	return matches[1], nil
}

// ExtractCharacterID extracts character ID from filename
func ExtractCharacterID(filename string) (string, error) {
	matches := characterIDRegex.FindStringSubmatch(filename)
	if len(matches) < 2 {
		return "", fmt.Errorf("invalid character filename format: %s", filename)
	}
	return matches[1], nil
}

// TryExtractUserName attempts to extract user name from .dat file
// Returns empty string if extraction fails (file is encrypted)
func TryExtractUserName(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// EVE .dat files are encrypted, but we can try to find readable strings
	// Look for common patterns that might contain username
	// This is a best-effort attempt and may not work for encrypted files

	content := string(data)

	// Try to find readable text that might be a username
	// Look for strings that look like usernames (alphanumeric, reasonable length)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and binary data
		if len(line) == 0 || len(line) > 100 {
			continue
		}
		// Check if line contains mostly printable ASCII
		if isPrintableASCII(line) && len(line) >= 3 && len(line) <= 50 {
			// Might be a username, but we can't be sure
			// Return empty to be safe since files are encrypted
		}
	}

	// Files are encrypted, so we can't reliably extract names
	return "", nil
}

// TryExtractCharacterName attempts to extract character name from .dat file
// Returns empty string if extraction fails (file is encrypted)
func TryExtractCharacterName(filePath string) (string, error) {
	// Same approach as TryExtractUserName
	// Since files are encrypted, we return empty string
	_, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Files are encrypted, so we can't reliably extract names
	return "", nil
}

// isPrintableASCII checks if a string contains only printable ASCII characters
func isPrintableASCII(s string) bool {
	for _, r := range s {
		if r < 32 || r > 126 {
			return false
		}
	}
	return true
}
