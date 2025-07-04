# Envelope

Envelope is a Terminal User Interface (TUI) for creating `.env` files by fetching secrets from your 1Password vaults.

## Prerequisites

Before using Envelope, you need to have the [1Password CLI](https://developer.1password.com/docs/cli/get-started/) installed and configured on your system. Envelope uses the `op` command-line tool to interact with your 1Password account.

## Installation

### Method 1: Download Pre-built Binary (Recommended)

#### macOS

1. Download the latest binary from the [releases page](https://github.com/gaganref/envelope/releases):

   - **Intel Macs:** `envelope-darwin-amd64`
   - **Apple Silicon Macs:** `envelope-darwin-arm64`

2. Make the binary executable (macOS only):

   ```bash
   chmod +x envelope-darwin-*
   ```

3. Run the application:
   ```bash
   ./envelope-darwin-amd64    # Intel Macs
   ./envelope-darwin-arm64    # Apple Silicon Macs
   ```

#### Windows

1. Download the latest binary from the [releases page](https://github.com/gaganref/envelope/releases):

   - **Windows:** `envelope-windows-amd64.exe`

2. Run the application:
   ```cmd
   ./envelope-windows-amd64.exe
   ```

### Method 2: Clone Repository and Build with Go

This method requires Go to be installed on your system.

#### macOS

```bash
git clone https://github.com/gaganref/envelope.git
cd envelope
go build .
./envelope
```

#### Windows

```cmd
git clone https://github.com/gaganref/envelope.git
cd envelope
go build .
envelope.exe
```

Alternatively, you can install directly using Go:

```bash
go install github.com/gaganref/envelope@latest
```

## Usage

### If you downloaded the binary:

**macOS:**

```bash
./envelope-darwin-amd64    # Intel Macs
./envelope-darwin-arm64    # Apple Silicon Macs
```

**Windows:**

```cmd
envelope-windows-amd64.exe
```

### If you built from source:

**macOS:**

```bash
./envelope
```

**Windows:**

```cmd
envelope.exe
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
   git commit
   git push
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. The GitHub Actions workflow will automatically:
   - Build binaries for macOS (Intel & Apple Silicon) and Windows
   - Create a GitHub release with the binaries attached
   - Generate release notes
