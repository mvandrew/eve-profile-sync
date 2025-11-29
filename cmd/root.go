package cmd

import (
	"fmt"
	"os"

	"eve-profile-sync/internal/backup"
	"eve-profile-sync/internal/config"
	"eve-profile-sync/internal/profile"
	"eve-profile-sync/internal/sync"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "eve-profile-sync",
	Short: "Synchronize EVE Online profile settings",
	Long: `EVE Profile Sync is a tool to synchronize EVE Online profile settings
by replacing all user and character configuration files with selected ones.`,
	Run: runSync,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runSync(cmd *cobra.Command, args []string) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Warning: Failed to load config: %v\n", err)
		cfg = &config.Config{}
	}

	// Step 1: Discover profiles directory
	profilesDir, err := discoverProfilesDirectory(cfg.ProfilesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Step 2: Select profile
	selectedProfile, err := selectProfile(profilesDir, cfg.Profile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Step 3: Select user file
	userFiles, err := profile.ListUserFiles(selectedProfile.Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to list user files: %v\n", err)
		os.Exit(1)
	}

	if len(userFiles) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No user files found in profile\n")
		os.Exit(1)
	}

	selectedUserFile, err := selectUserFile(userFiles, cfg.UserID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Step 4: Select character file
	charFiles, err := profile.ListCharacterFiles(selectedProfile.Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to list character files: %v\n", err)
		os.Exit(1)
	}

	if len(charFiles) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No character files found in profile\n")
		os.Exit(1)
	}

	selectedCharFile, err := selectCharacterFile(charFiles, cfg.CharacterID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Step 5: Show summary and confirm
	if !confirmOperation(selectedProfile, selectedUserFile, selectedCharFile) {
		fmt.Println("Operation cancelled.")
		os.Exit(0)
	}

	// Step 6: Validate operation
	if err := sync.ValidateOperation(selectedProfile.Path, selectedUserFile.Path, selectedCharFile.Path); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Validation failed: %v\n", err)
		os.Exit(1)
	}

	// Step 7: Create backup
	fmt.Println("Creating backup...")
	backupPath, err := backup.CreateBackup(selectedProfile.Path, selectedProfile.Name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create backup: %v\n", err)
		os.Exit(1)
	}

	// Verify backup
	if err := backup.VerifyBackup(backupPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Backup verification failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Backup created successfully: %s\n", backupPath)

	// Step 8: Perform synchronization
	fmt.Println("Synchronizing user files...")
	if err := sync.ReplaceUserFiles(selectedProfile.Path, selectedUserFile.Path); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to replace user files: %v\n", err)
		fmt.Printf("You can restore from backup: %s\n", backupPath)
		os.Exit(1)
	}

	fmt.Println("Synchronizing character files...")
	if err := sync.ReplaceCharacterFiles(selectedProfile.Path, selectedCharFile.Path); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to replace character files: %v\n", err)
		fmt.Printf("You can restore from backup: %s\n", backupPath)
		os.Exit(1)
	}

	// Step 9: Save configuration
	cfg.ProfilesDir = profilesDir
	cfg.Profile = selectedProfile.Name
	cfg.UserID = selectedUserFile.ID
	cfg.CharacterID = selectedCharFile.ID

	if err := config.SaveConfig(cfg); err != nil {
		fmt.Printf("Warning: Failed to save configuration: %v\n", err)
	}

	fmt.Println("Synchronization completed successfully!")
}

func discoverProfilesDirectory(savedDir string) (string, error) {
	// Try saved directory first
	if savedDir != "" {
		if err := profile.ValidateProfilesDirectory(savedDir); err == nil {
			return savedDir, nil
		}
	}

	// Try to find default directory
	profilesDir, err := profile.FindProfilesDirectory()
	if err == nil {
		return profilesDir, nil
	}

	// Prompt user for directory
	var userDir string
	prompt := &survey.Input{
		Message: "EVE profiles directory not found. Please enter the path:",
		Help:    "Default location: C:\\Users\\{user}\\AppData\\Local\\CCP\\EVE\\c_ccp_eve_online_tq_tranquility",
	}

	if err := survey.AskOne(prompt, &userDir, survey.WithValidator(survey.Required)); err != nil {
		return "", fmt.Errorf("failed to get directory from user: %w", err)
	}

	// Validate user-provided directory
	if err := profile.ValidateProfilesDirectory(userDir); err != nil {
		return "", fmt.Errorf("invalid directory: %w", err)
	}

	return userDir, nil
}

func selectProfile(profilesDir, savedProfile string) (*profile.Profile, error) {
	profiles, err := profile.ListProfiles(profilesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to list profiles: %w", err)
	}

	if len(profiles) == 0 {
		return nil, fmt.Errorf("no profiles found in directory: %s", profilesDir)
	}

	// Build options for survey
	options := make([]string, len(profiles))
	defaultIndex := 0
	for i, p := range profiles {
		options[i] = p.Name
		if savedProfile != "" && p.Name == savedProfile {
			defaultIndex = i
		}
	}

	var selected string
	prompt := &survey.Select{
		Message: "Select profile:",
		Options: options,
		Default: options[defaultIndex],
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		return nil, fmt.Errorf("failed to select profile: %w", err)
	}

	// Find selected profile
	for _, p := range profiles {
		if p.Name == selected {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("selected profile not found")
}

func selectUserFile(userFiles []profile.UserFile, savedUserID string) (*profile.UserFile, error) {
	// Build options for survey
	options := make([]string, len(userFiles))
	descriptions := make([]string, len(userFiles))
	defaultIndex := 0

	for i, uf := range userFiles {
		display := fmt.Sprintf("User ID: %s", uf.ID)
		if uf.Name != "" {
			display = fmt.Sprintf("%s (Name: %s)", display, uf.Name)
		}
		options[i] = display
		descriptions[i] = uf.Path

		if savedUserID != "" && uf.ID == savedUserID {
			defaultIndex = i
		}
	}

	var selected string
	prompt := &survey.Select{
		Message: "Select user file:",
		Options: options,
		Default: options[defaultIndex],
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		return nil, fmt.Errorf("failed to select user file: %w", err)
	}

	// Find selected user file
	for i, opt := range options {
		if opt == selected {
			return &userFiles[i], nil
		}
	}

	return nil, fmt.Errorf("selected user file not found")
}

func selectCharacterFile(charFiles []profile.CharacterFile, savedCharID string) (*profile.CharacterFile, error) {
	// Build options for survey
	options := make([]string, len(charFiles))
	defaultIndex := 0

	for i, cf := range charFiles {
		display := fmt.Sprintf("Character ID: %s", cf.ID)
		if cf.Name != "" {
			display = fmt.Sprintf("%s (Name: %s)", display, cf.Name)
		}
		options[i] = display

		if savedCharID != "" && cf.ID == savedCharID {
			defaultIndex = i
		}
	}

	var selected string
	prompt := &survey.Select{
		Message: "Select character file:",
		Options: options,
		Default: options[defaultIndex],
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		return nil, fmt.Errorf("failed to select character file: %w", err)
	}

	// Find selected character file
	for i, opt := range options {
		if opt == selected {
			return &charFiles[i], nil
		}
	}

	return nil, fmt.Errorf("selected character file not found")
}

func confirmOperation(selectedProfile *profile.Profile, selectedUserFile *profile.UserFile, selectedCharFile *profile.CharacterFile) bool {
	summary := fmt.Sprintf(`Operation Summary:
  Profile: %s
  User ID: %s
  Character ID: %s

This will replace all user and character files in the profile with the selected ones.
A backup will be created before making any changes.

Proceed?`, selectedProfile.Name, selectedUserFile.ID, selectedCharFile.ID)

	var proceed bool
	prompt := &survey.Confirm{
		Message: summary,
		Default: false,
	}

	if err := survey.AskOne(prompt, &proceed); err != nil {
		return false
	}

	return proceed
}
