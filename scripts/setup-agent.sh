#!/usr/bin/env bash
# GoPints Agent Installer
# Downloads the agent binary from GitHub Releases, installs it to
# /usr/local/bin, creates a config template, and registers a systemd service.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/rickcern44/gopints/main/scripts/setup-agent.sh | sudo bash
#   sudo bash setup-agent.sh [--version gopints-v1.2.3] [--server 192.168.1.10:9876]
#
# --version accepts either the full release tag (gopints-v1.2.3) or a bare
# version (v1.2.3) — the gopints- package prefix is added automatically if missing.

set -euo pipefail

# ── Configurable defaults ────────────────────────────────────────────────────
REPO="rickcern44/gopints"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/gopints"
CONFIG_FILE="${CONFIG_DIR}/config.json"
SERVICE_FILE="/etc/systemd/system/gopints-agent.service"
SERVICE_USER="gopints"
VERSION=""
SERVER_ADDR=""

# ── Argument parsing ─────────────────────────────────────────────────────────
while [[ $# -gt 0 ]]; do
  case "$1" in
    --version)  VERSION="$2";     shift 2 ;;
    --server)   SERVER_ADDR="$2"; shift 2 ;;
    *)          echo "Unknown argument: $1"; exit 1 ;;
  esac
done

# ── Helpers ──────────────────────────────────────────────────────────────────
info()    { echo "[gopints] $*"; }
success() { echo "[gopints] ✓ $*"; }
error()   { echo "[gopints] ERROR: $*" >&2; exit 1; }

need_root() {
  [[ $EUID -eq 0 ]] || error "This script must be run as root (use sudo)."
}

need_cmd() {
  command -v "$1" &>/dev/null || error "Required command not found: $1"
}

# ── Architecture detection ───────────────────────────────────────────────────
detect_arch() {
  local machine
  machine="$(uname -m)"
  case "$machine" in
    aarch64|arm64) echo "linux-arm64" ;;
    armv7l|armv7)  echo "linux-armv7" ;;
    x86_64)        echo "linux-amd64" ;;
    *) error "Unsupported architecture: $machine" ;;
  esac
}

# ── Resolve latest gopints release tag via GitHub API ────────────────────────
# The repo now produces two independent tag streams (gopints-vX.Y.Z, web-vX.Y.Z),
# so GitHub's single repo-wide "latest release" endpoint can't be used directly —
# it might point at a web-only release with no agent binary attached. List
# releases instead and take the newest one tagged gopints-v*.
resolve_version() {
  local latest
  latest="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases" \
    | grep '"tag_name"' | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/' \
    | grep '^gopints-v' | head -1)"
  [[ -n "$latest" ]] || error "Could not resolve latest gopints release from GitHub API."
  echo "$latest"
}

# ── Main ─────────────────────────────────────────────────────────────────────
need_root
need_cmd curl
need_cmd systemctl

ARCH="$(detect_arch)"
info "Detected platform: $ARCH"

if [[ -z "$VERSION" ]]; then
  info "Resolving latest gopints release…"
  VERSION="$(resolve_version)"
elif [[ "$VERSION" != gopints-* ]]; then
  VERSION="gopints-${VERSION}"
fi
info "Installing version: $VERSION"

BINARY_NAME="gopints-agent-${ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}"
INSTALL_PATH="${INSTALL_DIR}/gopints-agent"

# ── Download binary ───────────────────────────────────────────────────────────
info "Downloading ${DOWNLOAD_URL}…"
curl -fsSL --output "${INSTALL_PATH}.tmp" "$DOWNLOAD_URL"
chmod +x "${INSTALL_PATH}.tmp"
mv "${INSTALL_PATH}.tmp" "$INSTALL_PATH"
success "Installed binary to ${INSTALL_PATH}"

# ── Create service user (no login shell) ─────────────────────────────────────
if ! id -u "$SERVICE_USER" &>/dev/null; then
  useradd --system --no-create-home --shell /usr/sbin/nologin "$SERVICE_USER"
  success "Created system user: $SERVICE_USER"
fi

# ── Config file ───────────────────────────────────────────────────────────────
mkdir -p "$CONFIG_DIR"

if [[ ! -f "$CONFIG_FILE" ]]; then
  SERVER_PLACEHOLDER="${SERVER_ADDR:-"YOUR_SERVER_IP:9876"}"
  cat > "$CONFIG_FILE" <<EOF
{
  "server_addr": "${SERVER_PLACEHOLDER}",
  "taps": [
    {
      "id": 1,
      "gpio_pin": 17
    }
  ],
  "pulses_per_liter": 450,
  "flow_timeout_ms": 2000
}
EOF
  chown root:"$SERVICE_USER" "$CONFIG_FILE"
  chmod 640 "$CONFIG_FILE"
  success "Created config template at ${CONFIG_FILE}"
else
  info "Config already exists at ${CONFIG_FILE} — skipping."
fi

# ── systemd service ───────────────────────────────────────────────────────────
cat > "$SERVICE_FILE" <<EOF
[Unit]
Description=GoPints Kegerator Agent
Documentation=https://github.com/${REPO}
After=network.target
Wants=network.target

[Service]
Type=simple
User=${SERVICE_USER}
ExecStart=${INSTALL_PATH} --config ${CONFIG_FILE}
Restart=on-failure
RestartSec=5s

# Allow GPIO access
SupplementaryGroups=gpio
AmbientCapabilities=
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadOnlyPaths=/
ReadWritePaths=/tmp

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable gopints-agent
success "Registered and enabled gopints-agent.service"

# ── Start (or restart if already running) ────────────────────────────────────
if systemctl is-active --quiet gopints-agent; then
  systemctl restart gopints-agent
  success "Restarted gopints-agent"
else
  systemctl start gopints-agent
  success "Started gopints-agent"
fi

# ── Next steps ────────────────────────────────────────────────────────────────
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo " GoPints agent installed successfully!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
if [[ -z "$SERVER_ADDR" ]]; then
  echo " IMPORTANT: Edit the config before the agent can connect:"
  echo "   sudo nano ${CONFIG_FILE}"
  echo "   Set \"server_addr\" to your GoPints server IP and port (e.g. 192.168.1.10:9876)"
  echo "   Then restart: sudo systemctl restart gopints-agent"
  echo ""
fi
echo " Useful commands:"
echo "   sudo systemctl status  gopints-agent   # check status"
echo "   sudo journalctl -u     gopints-agent -f # stream logs"
echo "   sudo systemctl restart gopints-agent   # restart after config change"
echo "   sudo systemctl stop    gopints-agent   # stop the agent"
echo ""
