# Envelope

Envelope is a Terminal User Interface (TUI) for creating `.env` files by fetching secrets from your 1Password vaults.

## Prerequisites

Before using Envelope, you need to have the [1Password CLI](https://developer.1password.com/docs/cli/get-started/) installed and configured on your system. Envelope uses the `op` command-line tool to interact with your 1Password account.

## Usage

To use Envelope, simply run the application from your terminal:

```bash
go run .
```

The application will guide you through the following steps:

1.  **Enter Vault Name:** You'll be prompted to enter the name of the 1Password vault you want to access.
2.  **Select an Item:** A list of items from the specified vault will be displayed. You can navigate this list using the arrow keys and select an item by pressing `Enter`.
3.  **Enter Filename:** After selecting an item, you'll be prompted to enter a filename for your new `.env` file. It defaults to `.env` if you don't provide one.

Once completed, the `.env` file will be created in the current directory with the secrets from the selected 1Password item.
