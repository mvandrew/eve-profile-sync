# EVE Profile Sync

A command-line tool for synchronizing EVE Online profile settings by replacing all user and character configuration files within a profile with selected source files.

## Overview

EVE Profile Sync automates the process of synchronizing EVE Online profile settings across multiple user and character accounts. The tool discovers EVE profile directories, allows interactive selection of a source profile and specific user and character files, creates automatic backups, and replaces all matching configuration files within the selected profile.

The application runs as an interactive console program on Windows, guiding users through profile discovery, file selection, and synchronization operations. Selected preferences are saved to a configuration file for subsequent runs.

## Requirements

- Go 1.24 or later
- Windows 10/11
- EVE Online installed (for access to profile directories)

## Building

Build the executable using the provided batch script:

```bash
build.bat
```

Alternatively, build directly with Go:

```bash
go build -o eve-profile-sync.exe
```

The compiled binary `eve-profile-sync.exe` will be created in the project directory.

## Usage

Run the executable to start the interactive workflow:

```bash
eve-profile-sync.exe
```

The tool performs the following steps:

1. **Profile Directory Discovery**: Automatically detects the EVE profiles directory at `C:\Users\{user}\AppData\Local\CCP\EVE\c_ccp_eve_online_tq_tranquility`. If not found, prompts for the directory path.

2. **Profile Selection**: Lists all available profiles (directories matching `settings_*`) and prompts for selection. Previously selected profile is used as default.

3. **User File Selection**: Lists all user configuration files (`core_user_*.dat`) in the selected profile, displaying user IDs extracted from filenames. Previously selected user ID is used as default.

4. **Character File Selection**: Lists all character configuration files (`core_char_*.dat`) in the selected profile, displaying character IDs extracted from filenames. Previously selected character ID is used as default.

5. **Confirmation**: Displays a summary of the selected profile, user ID, and character ID, and requests confirmation before proceeding.

6. **Backup Creation**: Creates a timestamped ZIP backup of the entire profile directory in the `backup/` folder before making any modifications.

7. **Synchronization**: Replaces all `core_user_*.dat` files in the profile with the content of the selected user file, and replaces all `core_char_*.dat` files with the content of the selected character file. Original filenames are preserved.

8. **Configuration Save**: Saves the selected profile, user ID, and character ID to `config.yaml` for use as defaults in future runs.

## Configuration

The tool maintains a `config.yaml` file in the project directory with the following structure:

```yaml
profiles_dir: C:\Users\{user}\AppData\Local\CCP\EVE\c_ccp_eve_online_tq_tranquility
profile: ProfileName
user_id: "12345678"
character_id: "9876543210"
```

On first run, the configuration file is created automatically. Subsequent runs use saved values as defaults for interactive prompts. The configuration is updated after each successful synchronization.

## Backup

Backups are automatically created in the `backup/` directory before any file modifications. Backup files are named using the format:

```
settings_{ProfileName}_YYYYMMDD-HHmm.zip
```

For example: `settings_PVESolo_20251129-1745.zip`

Each backup contains a complete copy of the profile directory at the time of synchronization. Backup integrity is verified before proceeding with file replacements. If synchronization fails, the backup location is displayed for manual restoration.

## Project Structure

```
eve-profile-sync/
├── cmd/
│   └── root.go              # CLI command implementation and workflow orchestration
├── internal/
│   ├── profile/
│   │   ├── discover.go      # Profile directory discovery and validation
│   │   ├── parser.go         # File ID extraction from filenames
│   │   └── selector.go       # User and character file listing
│   ├── sync/
│   │   ├── replacer.go      # File replacement operations
│   │   └── validator.go     # Operation validation and safety checks
│   ├── backup/
│   │   └── creator.go        # ZIP backup creation and verification
│   └── config/
│       └── manager.go        # Configuration file management
├── backup/                   # Backup directory (created at runtime)
├── config.yaml               # Saved user preferences
├── main.go
└── go.mod
```

## Dependencies

- **github.com/spf13/cobra**: CLI framework for command structure and execution
- **github.com/spf13/viper**: Configuration management with YAML support
- **github.com/AlecAivazis/survey/v2**: Interactive terminal prompts for user input

