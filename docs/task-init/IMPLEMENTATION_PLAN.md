# Technical Specification and Implementation Plan

## Project Overview

**Project Name:** EVE Profile Sync

**Language:** Go 1.24+

**Platform:** Windows 10/11

**Type:** Console CLI Application

**Purpose:** Synchronize EVE Online profile settings by replacing all user and character configuration files with selected ones within a profile.

## Technical Requirements

### Functional Requirements

1. **Profile Directory Discovery**

   - Auto-detect EVE profiles directory: `C:\Users\{user}\AppData\Local\CCP\EVE\c_ccp_eve_online_tq_tranquility`
   - Fallback to user-provided path if default not found
   - Validate directory exists and is accessible

2. **Profile Selection**

   - List all profiles (directories matching pattern `settings_*`)
   - Interactive selection with profile name display
   - Support default selection from saved config

3. **User File Selection**

   - List all user files (matching pattern `core_user_*.dat`)
   - Display user ID extracted from filename
   - Attempt to extract user login from file (may not be possible due to encryption)
   - Support default selection from saved config

4. **Character File Selection**

   - List all character files (matching pattern `core_char_*.dat`)
   - Display character ID extracted from filename
   - Attempt to extract character name from file (may not be possible due to encryption)
   - Optionally filter by selected user (if relationship can be determined)
   - Support default selection from saved config

5. **Operation Summary**

   - Display selected profile, user, and character
   - Request user confirmation before proceeding

6. **Backup Creation**

   - Create full backup of profile directory
   - Backup location: `./backup/settings_{Profile}_YYYYMMDD-HHmm.zip`
   - Verify backup creation before proceeding with modifications

7. **File Synchronization**

   - Replace all `core_user_*.dat` files with selected user file content
   - Replace all `core_char_*.dat` files with selected character file content
   - Preserve original filenames (only content is replaced)

8. **Configuration Persistence**

   - Save selected profile, user ID, and character ID to config file
   - Load saved values as defaults on next run
   - Config format: YAML (using Viper)

### Non-Functional Requirements

- All user-facing messages in English
- Console-only interface (no GUI)
- Error handling with clear messages
- Validation of all inputs before operations
- Atomic operations where possible (backup before modify)

## Architecture Decisions

### Technology Stack

1. **CLI Framework:** `github.com/spf13/cobra` - Industry standard for Go CLI applications
2. **Configuration Management:** `github.com/spf13/viper` - Integrates with Cobra, supports YAML
3. **Interactive Prompts:** `github.com/AlecAivazis/survey` - Best for multi-step forms with validation
4. **File Operations:** Standard library (`os`, `path/filepath`, `io`, `archive/zip`)
5. **Error Handling:** Standard Go error patterns with custom error types

### Project Structure

```
eve-profile-sync/
├── cmd/
│   └── root.go              # Cobra root command
├── internal/
│   ├── profile/
│   │   ├── discover.go       # Profile directory discovery
│   │   ├── selector.go       # Profile/user/character selection
│   │   └── parser.go         # File ID extraction (attempt name extraction)
│   ├── sync/
│   │   ├── replacer.go       # File replacement operations
│   │   └── validator.go      # Operation validation
│   ├── backup/
│   │   └── creator.go         # Backup creation (zip)
│   └── config/
│       └── manager.go         # Config load/save (Viper)
├── backup/                   # Backup directory (created at runtime)
├── config.yaml               # Saved preferences
├── go.mod
├── go.sum
└── main.go
```

### Key Design Patterns

1. **Separation of Concerns:** Business logic in `internal/`, CLI in `cmd/`
2. **Dependency Injection:** Pass dependencies to functions for testability
3. **Error Wrapping:** Use `fmt.Errorf` with `%w` for error context
4. **Configuration First:** Load config early, use as defaults

## Component Breakdown

### 1. Profile Discovery (`internal/profile/discover.go`)

**Responsibilities:**

- Find EVE profiles directory (default path or user-provided)
- Validate directory existence and permissions
- List available profiles

**Key Functions:**

```go
func FindProfilesDirectory() (string, error)
func ListProfiles(profilesDir string) ([]Profile, error)
type Profile struct {
    Name string
    Path string
}
```

### 2. File Selector (`internal/profile/selector.go`)

**Responsibilities:**

- List user files in profile directory
- List character files in profile directory
- Extract IDs from filenames
- Attempt to extract names (may return empty if encrypted)

**Key Functions:**

```go
func ListUserFiles(profilePath string) ([]UserFile, error)
func ListCharacterFiles(profilePath string) ([]CharacterFile, error)
type UserFile struct {
    ID   string
    Name string // May be empty
    Path string
}
type CharacterFile struct {
    ID   string
    Name string // May be empty
    Path string
}
```

### 3. File Parser (`internal/profile/parser.go`)

**Responsibilities:**

- Extract user/character IDs from filenames
- Attempt to parse .dat files for names (graceful failure if encrypted)
- Handle file format variations

**Key Functions:**

```go
func ExtractUserID(filename string) (string, error)
func ExtractCharacterID(filename string) (string, error)
func TryExtractUserName(filePath string) (string, error) // May return empty
func TryExtractCharacterName(filePath string) (string, error) // May return empty
```

### 4. Backup Creator (`internal/backup/creator.go`)

**Responsibilities:**

- Create zip backup of profile directory
- Generate timestamped backup filename
- Verify backup integrity

**Key Functions:**

```go
func CreateBackup(profilePath, profileName string) (string, error)
func VerifyBackup(backupPath string) error
```

### 5. File Replacer (`internal/sync/replacer.go`)

**Responsibilities:**

- Replace all user files with selected user file
- Replace all character files with selected character file
- Preserve original filenames

**Key Functions:**

```go
func ReplaceUserFiles(profilePath, sourceUserFile string) error
func ReplaceCharacterFiles(profilePath, sourceCharFile string) error
```

### 6. Config Manager (`internal/config/manager.go`)

**Responsibilities:**

- Load saved configuration
- Save selected preferences
- Provide defaults for prompts

**Key Functions:**

```go
type Config struct {
    ProfilesDir string
    Profile     string
    UserID      string
    CharacterID string
}
func LoadConfig() (*Config, error)
func SaveConfig(cfg *Config) error
```

### 7. Root Command (`cmd/root.go`)

**Responsibilities:**

- Orchestrate the entire workflow
- Handle user interaction via Survey
- Coordinate all components
- Error handling and user feedback

**Workflow:**

1. Load config
2. Discover profiles directory
3. Select profile (with default)
4. Select user file (with default)
5. Select character file (with default)
6. Show summary and confirm
7. Create backup
8. Perform synchronization
9. Save config

## Implementation Steps

### Phase 1: Project Setup
1. Initialize Go module
2. Install dependencies (Cobra, Viper, Survey)
3. Create project structure
4. Set up Cobra root command

### Phase 2: Core Components
1. Implement profile discovery
2. Implement file listing and ID extraction
3. Implement config manager
4. Add basic error handling

### Phase 3: Interactive Interface
1. Integrate Survey for prompts
2. Implement profile selection
3. Implement user/character selection
4. Add summary and confirmation

### Phase 4: Backup and Sync
1. Implement backup creation (zip)
2. Implement file replacement operations
3. Add validation and error recovery
4. Test with sample data

### Phase 5: Polish and Testing
1. Add comprehensive error messages
2. Validate all edge cases
3. Test on Windows 10/11
4. Document usage

## Error Handling Strategy

1. **Validation Errors:** Check inputs before operations, return clear errors
2. **File System Errors:** Check permissions, disk space, file locks
3. **Backup Errors:** Verify backup before proceeding, abort if backup fails
4. **Sync Errors:** Log errors, attempt rollback if possible
5. **User Errors:** Provide helpful messages, suggest solutions

## Testing Considerations

1. **Unit Tests:** Test each component in isolation
2. **Integration Tests:** Test full workflow with mock data
3. **Error Cases:** Test all error paths
4. **Windows Compatibility:** Test on actual Windows environment

## Dependencies

```go
require (
    github.com/spf13/cobra v1.8.0
    github.com/spf13/viper v1.18.0
    github.com/AlecAivazis/survey/v2 v2.3.7
)
```

## Notes

- EVE Online .dat files are encrypted, so extracting user/character names may not be possible
- File operations should preserve original filenames (only content is replaced)
- Backup is critical - always verify before proceeding
- Config file should be in project directory for portability

