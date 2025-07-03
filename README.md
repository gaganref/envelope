# Envelope

Envelope is a Terminal User Interface (TUI) for creating `.env` files by fetching secrets from your 1Password vaults.

## Prerequisites

Before using Envelope, you need to have the [1Password CLI](https://developer.1password.com/docs/cli/get-started/) installed and configured on your system. Envelope uses the `op` command-line tool to interact with your 1Password account.

## Installation

### Download Pre-built Binary

Download the latest binary for your platform from the [releases page](https://github.com/gaganref/envelope/releases):

- **macOS (Intel):** `envelope-darwin-amd64`
- **macOS (Apple Silicon):** `envelope-darwin-arm64`
- **Windows:** `envelope-windows-amd64.exe`

Make the binary executable (macOS/Linux only):

```bash
chmod +x envelope-darwin-*
```

### Build from Source

If you have Go installed, you can build from source:

```bash
go install github.com/gaganref/envelope@latest
```

Or clone and build locally:

```bash
git clone https://github.com/gaganref/envelope.git
cd envelope
go build .
```

## Usage

Run the application from your terminal:

```bash
./envelope-darwin-amd64  # macOS Intel
./envelope-darwin-arm64  # macOS Apple Silicon
./envelope-windows-amd64.exe  # Windows
```

Or if built from source:

```bash
./envelope
```

The application will guide you through the following steps:

1.  **Enter Vault Name:** You'll be prompted to enter the name of the 1Password vault you want to access.
2.  **Select an Item:** A list of items from the specified vault will be displayed. You can navigate this list using the arrow keys and select an item by pressing `Enter`.
3.  **Enter Filename:** After selecting an item, you'll be prompted to enter a filename for your new `.env` file. It defaults to `.env` if you don't provide one.

Once completed, the `.env` file will be created in the current directory with the secrets from the selected 1Password item.

## Creating a Release

To create a new release with binaries:

1. Tag your commit with a version:

   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. The GitHub Actions workflow will automatically:
   - Build binaries for macOS (Intel & Apple Silicon) and Windows
   - Create a GitHub release with the binaries attached
   - Generate release notes
