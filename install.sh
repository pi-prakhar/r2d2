#!/bin/bash

# Detect the platform (macOS, Linux, or Windows)
OS=$(uname -s)
ARCH=$(uname -m)

echo "Installing r2d2..."
echo "OS: $OS"
echo "ARCH: $ARCH"

# GitHub repository details
REPO="pi-prakhar/r2d2"
GITHUB_API="https://api.github.com/repos/$REPO/releases/latest"

# Fetch the latest release information
LATEST_RELEASE=$(curl -s $GITHUB_API)

# Parse the latest release version and assets
LATEST_VERSION=$(echo $LATEST_RELEASE | jq -r '.tag_name')
ASSET_NAME="r2d2-"

echo "LATEST_VERSION: $LATEST_VERSION"


# Determine appropriate asset for the user's platform and architecture
if [[ "$OS" == "Darwin" ]]; then
  if [[ "$ARCH" == "x86_64" ]]; then
    ASSET_NAME+="darwin-amd64.tar.gz"
  elif [[ "$ARCH" == "arm64" ]]; then
    ASSET_NAME+="darwin-arm64.tar.gz"
  fi
elif [[ "$OS" == "Linux" ]]; then
  if [[ "$ARCH" == "x86_64" ]]; then
    ASSET_NAME+="linux-amd64.tar.gz"
  fi
fi

echo "ASSET_NAME: $ASSET_NAME"

# Check if the asset exists for the platform
if [[ "$ASSET_NAME" == "r2d2-"* ]]; then
  echo "Downloading $ASSET_NAME for $OS ($ARCH)..."
else
  echo "Unsupported platform or architecture."
  exit 1
fi

# Download the asset
curl -LO "https://github.com/$REPO/releases/download/$LATEST_VERSION/$ASSET_NAME"

# Extract the binary
if [[ "$OS" == "Darwin" || "$OS" == "Linux" ]]; then
  tar -xzf "$ASSET_NAME"
  chmod +x r2d2

  # macOS specific: remove quarantine attribute
  if [[ "$OS" == "Darwin" ]]; then
    xattr -d com.apple.quarantine ./r2d2 || true
  fi

  # Move to /usr/local/bin
  sudo mv ./r2d2 /usr/local/bin/
  
  echo "r2d2 binary installed successfully."
else
  echo "Installation failed. Unsupported OS."
  exit 1
fi
