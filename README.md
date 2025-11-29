# EVE Online Multiboxing Profile Sync — UI Settings Synchronizer

Multiboxers often waste time manually copying `core_user_*.dat` and `core_char_*.dat` files between launcher profiles to keep UI layouts, overview, and hotkeys consistent across accounts. **EVE Profile Sync** automates this by syncing your chosen user and character settings across all files in a profile. Ideal for players running multiple EVE Online accounts on the same PC.

---

## Why Multiboxers Need This

Keeping UI settings synchronized between profiles is a well-known pain:

- Each launcher profile creates separate EVE settings directories  
- UI, window layout, overview, and hotkeys reset per profile  
- After reinstalling EVE or switching PCs, configuration becomes inconsistent  
- Multiboxers (2–20+ clients) lose hours manually copying files  

**This tool fixes that problem** by letting you pick a "source" UI and push it instantly to every user/character settings file inside a launcher profile.

---

## Features

- **Sync UI settings across all accounts** by replacing all matching `core_user_*.dat` and `core_char_*.dat`  
- **Preserves filenames** while overwriting file contents  
- **Automatic launcher profile discovery**:  
  `C:\Users\{user}\AppData\Local\CCP\EVE\c_ccp_eve_online_tq_tranquility`
- **Interactive CLI workflow** with profile, user, and character selectors  
- **Timestamped ZIP backups** before any changes  
- **Persistent configuration (`config.yaml`)** remembering previously used values  
- **Windows-ready executable** — runs on Windows 10/11  
- Designed specifically for **EVE Online multiboxers**  

---

## How It Works

The tool performs the following steps:

1. **Profile Directory Discovery**: Automatically detects the EVE profiles directory at `C:\Users\{user}\AppData\Local\CCP\EVE\c_ccp_eve_online_tq_tranquility`. If not found, prompts for the directory path.

2. **Profile Selection**: Lists all available profiles (directories matching `settings_*`) and prompts for selection. Previously selected profile is used as default.

3. **User File Selection**: Lists all user configuration files (`core_user_*.dat`) in the selected profile, displaying user IDs extracted from filenames. Previously selected user ID is used as default.

4. **Character File Selection**: Lists all character configuration files (`core_char_*.dat`) in the selected profile, displaying character IDs extracted from filenames. Previously selected character ID is used as default.

5. **Confirmation**: Displays a summary of the selected profile, user ID, and character ID, and requests confirmation before proceeding.

6. **Backup Creation**: Creates a timestamped ZIP backup of the entire profile directory in the `backup/` folder before making any modifications.

7. **Synchronization**: Replaces all `core_user_*.dat` files in the profile with the content of the selected user file, and replaces all `core_char_*.dat` files with the content of the selected character file. Original filenames are preserved.

8. **Configuration Save**: Saves the selected profile, user ID, and character ID to `config.yaml` for use as defaults in future runs.

This guarantees that UI layout, overview settings, and hotkeys remain identical across all your accounts.

---

## Requirements

- Windows 10/11
- EVE Online installed (for access to profile directories)
- Pre-built binary or Go 1.24+ for building from source

---

## Installation

### Option 1 — Pre-built Release
Download `eve-profile-sync.exe` from the Releases page.

### Option 2 — Build from Source
Requires Go 1.24+:

```bash
go build -o eve-profile-sync.exe
```

Or use the included build script:

```bash
build.bat
```

The compiled binary `eve-profile-sync.exe` will be created in the project directory.

---

## Usage

Run the executable to start the interactive workflow:

```bash
eve-profile-sync.exe
```

The application runs as an interactive console program on Windows, guiding users through profile discovery, file selection, and synchronization operations. Selected preferences are saved to a configuration file for subsequent runs, making repeated synchronizations faster and more convenient.

Workflow:

1. Select launcher profile
2. Select source user (`core_user_*.dat`)
3. Select source character (`core_char_*.dat`)
4. Confirm sync
5. Backup is created automatically
6. UI settings are synchronized across all profiles

Everything happens inside a simple interactive console workflow.

---

## Configuration

The tool maintains a `config.yaml` file in the project directory with the following structure:

```yaml
profiles_dir: C:\Users\{user}\AppData\Local\CCP\EVE\c_ccp_eve_online_tq_tranquility
profile: ProfileName
user_id: "12345678"
character_id: "9876543210"
```

On first run, the configuration file is created automatically. Subsequent runs use saved values as defaults for interactive prompts. The configuration is updated after each successful synchronization.

---

## Backup Strategy

Before making changes, a ZIP archive is created in the `backup/` directory. Backup files are named using the format:

```
settings_{ProfileName}_YYYYMMDD-HHmm.zip
```

For example: `settings_PVESolo_20251129-1745.zip`

Each backup contains a complete copy of the profile directory at the time of synchronization. Backup integrity is verified before proceeding with file replacements. If synchronization fails, the backup location is displayed for manual restoration.

---

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

---

## Dependencies

- **github.com/spf13/cobra**: CLI framework for command structure and execution
- **github.com/spf13/viper**: Configuration management with YAML support
- **github.com/AlecAivazis/survey/v2**: Interactive terminal prompts for user input

---

## Roadmap

* Cross-profile synchronization
* Optional cloud backup & restore
* Overview/hotkeys diff checking
* Auto-detection of active clients

---

## License

Licensed under MIT. Contributions welcome.
