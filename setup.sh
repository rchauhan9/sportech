#!/usr/bin/env bash

set -euo pipefail

trap on_error ERR

PYTHON_VERSION="${PYTHON_VERSION:-3.11.0}"
GO_VERSION="${GO_VERSION:-1.18.3}"

BREW_TAPS=(
  "homebrew/core"
  "homebrew/cask"
)

BREW_FORMULAE=(
  "awscli"
  "cocoapods"
  "go"  # we use gvm to install go, but compiling go requires go, so the one from brew bootstraps it
  "golang-migrate"
  "helm"
  "jq"
  "kubectl"
  "kubectx"
  "pre-commit"
  "pyenv"
  "python3"
)

BREW_CASKS=(
  "docker"
)

PIP_PACKAGES=(
  "poetry"
)

echo_coloured() {
  echo -e "\033[0;$2m$1\033[0m"
}

info() {
  echo_coloured "${1:-}" 34  # blue
}

err() {
  echo_coloured "${1:-}" 31  # red
}

on_error() {
  err "Something went wrong with the install script."
}

install_single() {
  local install_cmd="$1"
  shift
  info "Installing with '$install_cmd'..."
  echo "$@" | xargs -n1 $install_cmd
}

install_multi() {
  local install_cmd="$1"
  shift
  info "Installing with '$install_cmd'..."
  $install_cmd "$@"
}

install_external() {
  local cmd="$1"
  local uri="$2"
  if ! command -v "$cmd" &>/dev/null; then
    info "Installing $cmd..."
    bash -c "$(curl -fsSL $uri)"
  else
    info "$cmd is already installed!"
  fi
}

info "Welcome to the Sports API setup script!"
info "This will setup some standard tools that we use for developing."
info

# Global gitignore

info "Setting up global gitignore..."
cat >"$HOME/.gitignore_global" <<EOL
.idea
.vscode
.iml
.gradle
EOL
git config --global core.excludesfile "$HOME/.gitignore_global"

# Generic package setup

info "Installing generic packages..."

if ! xcode-select --print-path &>/dev/null; then
  info "Installing Xcode command line tools..."
  xcode-select --install
else
  info "Xcode command line tools already installed!"
fi

install_external "brew" "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh" && brew update && brew upgrade
install_external "sdk" "https://get.sdkman.io?rcupdate=false"
if [[ ! -d "$HOME/.gvm" ]]; then
  GVM_NO_UPDATE_PROFILE=true install_external "gvm" "https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer"
else
  info "GVM is already installed!"
fi

install_single "brew tap" "${BREW_TAPS[@]}"

install_multi "brew install" "${BREW_FORMULAE[@]}"
install_multi "brew install --cask" "${BREW_CASKS[@]}"

# Python setup

info "Setting up python..."
pyenv install --skip-existing "${PYTHON_VERSION}"
pyenv global "${PYTHON_VERSION}"
install_multi "pip3 install" "${PIP_PACKAGES[@]}"

# Go setup

info "Setting up go..."
set +u  # gvm has some unbound variables (which it checks for in the script)
. "$HOME/.gvm/scripts/gvm"
gvm install "go${GO_VERSION}"
gvm use "go${GO_VERSION}" --default
set -u

# Pre-commit setup

info "Setting up pre-commit..."
pre-commit install -t pre-commit -t pre-push

# Docker setup

info "Setting up docker..."
# Add a default daemon.json
# This sets the default address pools as the range used by docker conflicts with some of our subnets
mkdir -p ~/.docker
touch ~/.docker/daemon.json
cat > ~/.docker/daemon.json <<EOF
{
  "default-address-pools": [
    {
      "base": "172.200.0.0/16",
      "size": 24
    }
  ],
  "experimental": false,
  "features": {
    "buildkit": true
  }
}
EOF
if ! command -v docker &>/dev/null; then
  open /Applications/Docker.app
  info "A popup will appear asking you to give privileges for docker - please allow it and the script will continue when the engine has started"
fi
while ! command -v docker &>/dev/null || ! docker ps &>/dev/null; do
  sleep 5
done

info "Setup completed successfully! 🎉"
