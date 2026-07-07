#!/usr/bin/env sh
# Moodwave CLI — Unix install script
# Usage: curl -sSL https://raw.githubusercontent.com/Boredooms/Moodwave-CLI/main/cli/scripts/install.sh | sh
#
# This script:
# 1. Detects OS and architecture
# 2. Downloads the appropriate binary from GitHub Releases
# 3. Installs it to /usr/local/bin or ~/bin (if not root)
# 4. Verifies the installation

set -e

REPO="Boredooms/Moodwave-CLI"
BINARY="moodwave"
VERSION="${MOODWAVE_VERSION:-latest}"
INSTALL_DIR=""

# ── Colors ────────────────────────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m'

info()    { printf "${BLUE}info${NC}  %s\n" "$1"; }
success() { printf "${GREEN}ok${NC}    %s\n" "$1"; }
warn()    { printf "${YELLOW}warn${NC}  %s\n" "$1"; }
error()   { printf "${RED}error${NC} %s\n" "$1" >&2; exit 1; }

# ── Detect OS ─────────────────────────────────────────────────────────────────
detect_os() {
  case "$(uname -s)" in
    Linux)  echo "linux" ;;
    Darwin) echo "darwin" ;;
    *)      error "Unsupported OS: $(uname -s). Please build from source." ;;
  esac
}

# ── Detect architecture ───────────────────────────────────────────────────────
detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64) echo "amd64" ;;
    arm64|aarch64) echo "arm64" ;;
    armv7l) echo "arm" ;;
    *)  error "Unsupported architecture: $(uname -m). Please build from source." ;;
  esac
}

# ── Find install directory ────────────────────────────────────────────────────
find_install_dir() {
  if [ -w "/usr/local/bin" ]; then
    echo "/usr/local/bin"
  elif [ -w "$HOME/.local/bin" ]; then
    mkdir -p "$HOME/.local/bin"
    echo "$HOME/.local/bin"
  elif [ -w "$HOME/bin" ]; then
    mkdir -p "$HOME/bin"
    echo "$HOME/bin"
  else
    mkdir -p "$HOME/.local/bin"
    echo "$HOME/.local/bin"
  fi
}

# ── Download with progress ────────────────────────────────────────────────────
download() {
  url="$1"
  dest="$2"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL --progress-bar "$url" -o "$dest"
  elif command -v wget >/dev/null 2>&1; then
    wget -q --show-progress "$url" -O "$dest"
  else
    error "Neither curl nor wget found. Please install one of them."
  fi
}

# ── Resolve version ───────────────────────────────────────────────────────────
resolve_version() {
  if [ "$VERSION" = "latest" ]; then
    if command -v curl >/dev/null 2>&1; then
      VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed 's/.*"tag_name": "\(.*\)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
      VERSION=$(wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed 's/.*"tag_name": "\(.*\)".*/\1/')
    fi
    if [ -z "$VERSION" ]; then
      error "Could not resolve latest version. Set MOODWAVE_VERSION=vX.Y.Z to install a specific version."
    fi
  fi
  echo "$VERSION"
}

# ── Main ──────────────────────────────────────────────────────────────────────
main() {
  printf "\n${BOLD}Moodwave CLI Installer${NC}\n"
  printf "%s\n\n" "$(printf '%.0s─' $(seq 1 40))"

  OS=$(detect_os)
  ARCH=$(detect_arch)
  VERSION=$(resolve_version)
  INSTALL_DIR=$(find_install_dir)

  info "OS:         $OS"
  info "Arch:       $ARCH"
  info "Version:    $VERSION"
  info "Install to: $INSTALL_DIR"
  printf "\n"

  # Build the download URL.
  BINARY_NAME="${BINARY}-${OS}-${ARCH}"
  DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}"

  # Create temp directory.
  TMP_DIR=$(mktemp -d)
  TMP_BINARY="$TMP_DIR/$BINARY"

  # Download.
  info "Downloading ${BINARY_NAME}..."
  if ! download "$DOWNLOAD_URL" "$TMP_BINARY"; then
    rm -rf "$TMP_DIR"
    error "Download failed from: $DOWNLOAD_URL"
  fi

  # Make executable.
  chmod +x "$TMP_BINARY"

  # Verify binary.
  if ! "$TMP_BINARY" --version >/dev/null 2>&1; then
    rm -rf "$TMP_DIR"
    error "Downloaded binary failed to run. Try building from source."
  fi

  # Install.
  INSTALL_PATH="$INSTALL_DIR/$BINARY"
  mv "$TMP_BINARY" "$INSTALL_PATH"
  rm -rf "$TMP_DIR"

  success "Installed to $INSTALL_PATH"

  # Check PATH.
  if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    warn "$INSTALL_DIR is not in your PATH."
    warn "Add this to your shell profile:"
    printf "  export PATH=\"\$PATH:$INSTALL_DIR\"\n"
  fi

  # Verify.
  if command -v "$BINARY" >/dev/null 2>&1; then
    INSTALLED_VERSION=$("$BINARY" --version 2>&1)
    success "Installed: $INSTALLED_VERSION"
  fi

  printf "\n${BOLD}Ready! Try:${NC}\n"
  printf "  moodwave init\n"
  printf "  moodwave scan\n"
  printf "  moodwave play\n\n"
}

main "$@"
