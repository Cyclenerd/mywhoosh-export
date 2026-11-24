# MyWhoosh FIT File Exporter

[![Badge: Linux](https://img.shields.io/badge/Linux-FCC624.svg?logo=linux&logoColor=black)](#readme)
[![Badge: macOS](https://img.shields.io/badge/macOS-000000.svg?logo=apple&logoColor=white)](#readme)
[![Badge: Windows](https://img.shields.io/badge/Windows-008080.svg?logoColor=white)](#readme)
[![Badge: Go](https://img.shields.io/badge/Go-%2300ADD8.svg?logo=go&logoColor=white)](#readme)
[![Badge: Vibe coded with OpenCode](https://img.shields.io/badge/Vibe%20coded%20with-OpenCode-000000.svg)](https://opencode.ai)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

A minimal, single-file Go CLI tool to automatically export and sync FIT files from MyWhoosh to your local machine. Downloads the last 10 activities from your MyWhoosh account.

## Features

- 🔐 Secure authentication via environment variables only
- 📥 Automatic discovery and download of last 10 FIT files from MyWhoosh
- 🔄 Smart sync - skips already downloaded files
- 📊 Activity listing and filtering
- ⚡ Simple sequential downloads
- 🚀 Zero external dependencies

## Getting Started

You can either download a pre-built binary for your operating system or build the application from source.

### Pre-built Binaries (Recommended for most users)

This project require you to use a command-line interface.
Don't worry, it's easier than it looks!

#### <img src="https://www.microsoft.com/favicon.ico" width=16 height=16/> Windows (PowerShell)

1. [Download](https://github.com/Cyclenerd/mywhoosh-export/releases/latest) the `mywhoosh-export.exe` application for Windows and your CPU architecture. If you are unsure, usually `x86_64` will work.
    * [x86_64](https://github.com/Cyclenerd/mywhoosh-export/releases/latest/download/mywhoosh-export-windows-x86_64.exe) Intel or AMD 64-Bit CPU
    * [arm64](https://github.com/Cyclenerd/mywhoosh-export/releases/latest/download/mywhoosh-export-windows-arm64.exe) Arm-based 64-Bit CPU
1. Rename it to: `mywhoosh-export.exe`
1. Press the <kbd>Windows key</kbd> + <kbd>X</kbd>.
1. Choose "Windows PowerShell".
1. Type `cd` followed by the path to your your folder where you saved the application (e.g., `cd C:\Users\YourName\Documents\MyWhooshExport`) and press Enter.
1. Type `.\mywhoosh-export.exe` and press Enter.

### <img src="https://developer.apple.com/favicon.ico" width=16 height=16/> macOS (Terminal)

1. [Download](https://github.com/Cyclenerd/mywhoosh-export/releases/latest) the `mywhoosh-export` application for macOS and your CPU architecture.
    * [x86_64](https://github.com/Cyclenerd/mywhoosh-export/releases/latest/download/mywhoosh-export-macos-x86_64) Intel 64-bit CPU
    * [arm64](https://github.com/Cyclenerd/mywhoosh-export/releases/latest/download/mywhoosh-export-macos-arm64) Apple silicon 64-bit CPU (M1, M2...)
1. Rename it to: `mywhoosh-export`
1. Press <kbd>Command</kbd> + <kbd>Space</kbd>.
1. Type "Terminal" and press Enter.
1. Type `cd` followed by the path to your folder where you saved the application (e.g., `cd /Users/YourName/Documents/MyWhooshExport`) and press Enter.
1. Type `chmod +x mywhoosh-export` and press Enter. (This makes the tool work).
1. Type `./mywhoosh-export` and press Enter.
1. (Optional) Allow unsigned app on macOS
    1. Go to System Settings
    1. Click on Privacy & Security
    1. Scroll down to Security section
    1. Look for a message about the blocked `mywhoosh-export` app and click "Open Anyway"

### <img src="https://www.kernel.org/theme/images/logos/favicon.png" width=16 height=16/> Linux (Terminal)

1. [Download](https://github.com/Cyclenerd/mywhoosh-export/releases/latest) the `mywhoosh-export` application for Linux and your CPU architecture.
    * [x86_64](https://github.com/Cyclenerd/mywhoosh-export/releases/latest/download/mywhoosh-export-linux-x86_64) Intel or AMD 64-Bit CPU
    * [arm64](https://github.com/Cyclenerd/mywhoosh-export/releases/latest/download/mywhoosh-export-linux-arm64) Arm-based 64-Bit CPU (i.e. Raspberry Pi)
1. Rename it to: `mywhoosh-export`
1. Open a Terminal.
1. Type `cd` followed by the path to your your folder where you saved the application (e.g., `cd /home/YourName/MyWhooshExport`) and press Enter.
1. Type `chmod +x mywhoosh-export` and press Enter. (This makes the tool work).
1. Type `./mywhoosh-export` and press Enter.

### Building from Source

If you're a developer or want to modify the code, you can build the application from source.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/Cyclenerd/mywhoosh-export.git
    cd mywhoosh-export
    ```
2.  **Build the application:**
    ```bash
    go build .
    ```
3.  **Run the application:**
    ```bash
    ./mywhoosh-export
    ```

## Getting Your Credentials

You need two pieces of information:
1. **User ID** - Your unique identifier in MyWhoosh
2. **Auth Token** - A token for API authentication

### How to Obtain Credentials

1. Log in to MyWhoosh (https://event.mywhoosh.com/auth/login)
2. Open the Web Developer Tools in your Browser: <kbd>F12</kbd> or <kbd>Ctrl+Shift+I</kbd>
3. Go to the Network tab
4. Navigate to your activities page
5. Find any API request to `service14.mywhoosh.com`
6. Look at the request cookies:
   - `whoosh_token` - This is your Auth Token
   - `whoosh_uuid` - This is your User ID

## Quick Start

Set your credentials as environment variables:

```bash
export MYWHOOSH_TOKEN="your-token-here"
export MYWHOOSH_USER_ID="your-user-id-here"
```

Then download your FIT files:

```bash
./mywhoosh-export -output ~/MyWhoosh-Sync
```

Done! Your FIT files are in `~/MyWhoosh-Sync/`

## Usage

```text
Usage of ./mywhoosh-export:
  -list
        List activities without downloading
  -output string
        Output directory for FIT files (default "./fit-files")
  -token string
        Authentication token (or MYWHOOSH_TOKEN env var)
  -user-id string
        User ID (or MYWHOOSH_USER_ID env var)
```

## Examples

### List all activities without downloading

```bash
export MYWHOOSH_TOKEN="your_token"
export MYWHOOSH_USER_ID="your_id"
./mywhoosh-export -list
```

### Download all FIT files

```bash
export MYWHOOSH_TOKEN="your_token"
export MYWHOOSH_USER_ID="your_id"
./mywhoosh-export -output ~/MyWhoosh-Sync
```

### Use command-line arguments

```bash
./mywhoosh-export -token "your_token" -user-id "your_id" -output ~/MyWhoosh-Sync
```

### Sync again (skips existing files)

```bash
./mywhoosh-export -output ~/MyWhoosh-Sync
```

## How It Works

1. **Authentication**: Uses token from environment or command-line flags
2. **Discovery**: Fetches activity list from MyWhoosh API (`POST /v2/rider/profile/activities`) - last 10 activities only
3. **Presigned URLs**: Requests presigned S3 URLs from MyWhoosh API (`POST /v2/rider/profile/download-activity-file`)
4. **Download**: Downloads each FIT file from S3 using presigned URL (valid for 60 seconds)
5. **Smart Sync**: Skips files that already exist locally
6. **Reporting**: Provides download statistics and error logging

## API Endpoints

The tool uses the following MyWhoosh API endpoints:

- `POST https://service14.mywhoosh.com/v2/rider/profile/activities` - List activities
  - Request: `{ "page": 1, "sortDate": "DESC" }`
  - Response: List of activities with file IDs

- `POST https://service14.mywhoosh.com/v2/rider/profile/download-activity-file` - Get presigned URL
  - Request: `{ "fileId": "<activityFileId>" }`
  - Response: Presigned AWS S3 URL

- `https://mywhooshprod.s3.eu-west-1.amazonaws.com` - S3 direct download with presigned URL

## File Naming

Downloaded FIT files follow this naming convention:

```
YYYYMMDD_HHMMSS_<fileId>.fit
```

Example: `20261120_161105_6b007f7467d421da7b053f50.fit`

- `YYYYMMDD` - Activity date
- `HHMMSS` - Activity start time
- `<fileId>` - Unique file identifier from MyWhoosh

## Troubleshooting

### Error: "API error: status 401"

Your authentication token has expired or is invalid. Update your environment variables with a new token.

### Error: "failed to fetch activities"

Check that your user ID and token are correct. Verify by looking at the Network tab in DevTools.

### Error: "S3 download failed: status 403"

This should not happen as presigned URLs are requested from the API. If it occurs, the presigned URL may have expired (valid for 60 seconds). Retry the download.

## License

This project is licensed under the **GNU Affero General Public License v3 (AGPL-3.0)**.

See [LICENSE](LICENSE) file for details.
